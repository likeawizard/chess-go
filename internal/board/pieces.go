package board

const PieceOffset = 6

var Pieces = [12]string{"P", "B", "N", "R", "Q", "K", "p", "b", "n", "r", "q", "k"}
var PiecesUnicode = [12]string{"\u265F", "\u265D", "\u265E", "\u265C", "\u265B", "\u265A", "\u2659", "\u2657", "\u2658", "\u2656", "\u2655", "\u2654"}
var Sqaures = [2]string{"\u2B1B", "\u2B1C"}

func PieceSymbolToInt(piece string) int {
	for i, p := range Pieces {
		if p == piece {
			return i + 1
		}
	}
	return 0
}
func GetKing(b Board, color string) (c Coord) {
	c.File = 8
	c.Rank = 8
	king := 6
	if color == BlackToMove {
		king += PieceOffset
	}
	for f, file := range b.Coords {
		for r, piece := range file {
			if piece == king {
				c.File = f
				c.Rank = r
			}
		}
	}
	return
}

func GetPieces(b Board, color string) (pieces []Coord) {
	for f, file := range b.Coords {
		for r := range file {
			coord := Coord{f, r}
			piece1, color1 := GetPiece(b, coord)
			if piece1 != 0 && color == color1 {
				pieces = append(pieces, coord)
			}
		}
	}
	return
}

func GetPiece(b Board, coord Coord) (piece int, color string) {
	piece = b.AccessCoord(coord)
	if piece <= PieceOffset {
		color = WhiteToMove
	} else {
		color = BlackToMove
	}
	return
}

func (b Board) MoveToPretty(move string) (pretty string) {
	from, to := longAlgToCoords(move)
	targetPiece := b.AccessCoord(to)
	piece := b.AccessCoord(from)
	switch {
	case piece == P || piece == p:
		pretty = move[2:]
		if move[:1] != move[2:3] {
			pretty = move[:1] + "x" + pretty
		}
	case move == CastlingMoves[0] || move == CastlingMoves[2]:
		return "O-O-O"
	case move == CastlingMoves[1] || move == CastlingMoves[3]:
		return "O-O"
	default:
		pretty = Pieces[(piece-1)%PieceOffset]
		if targetPiece > 0 {
			pretty += "x"
		}
		pretty += move[2:]
	}

	return
}

func (b Board) hasRankConflict(from Coord) bool {
	identicalPieceCount := 0
	piece := b.AccessCoord(from)
	for i := 0; i < 8; i++ {
		if b.Coords[from.File][i] == piece {
			identicalPieceCount++
		}
	}
	return identicalPieceCount > 1
}

func (b Board) hasFileConflict(from Coord) bool {
	identicalPieceCount := 0
	piece := b.AccessCoord(from)
	for i := 0; i < 8; i++ {
		if b.Coords[i][from.Rank] == piece {
			identicalPieceCount++
		}
	}
	return identicalPieceCount > 1
}

func (b Board) GetMoves(color string) (moves, captures []string) {
	return b.getMoves(color, false)
}

func (b Board) GetMovesNoCastling(color string) (moves, captures []string) {
	return b.getMoves(color, true)
}

func (b Board) getMoves(color string, excludeCastling bool) (moves, captures []string) {
	pieces := GetPieces(b, color)
	for _, piece := range pieces {
		m, c := b.GetAvailableMovesRaw(piece, excludeCastling)
		moves = append(moves, m...)
		captures = append(captures, c...)
	}
	return
}

func (b *Board) playOutLine(line []string) {
	for _, move := range line {
		b.MoveLongAlg(move)
	}
}

func (b *Board) placePiece(coord string, piece int) {
	target := AlgToCoord(coord)
	b.Coords[target.File][target.Rank] = piece
}
