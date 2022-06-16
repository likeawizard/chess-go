package board

const PieceOffset = 6

var Pieces = [12]string{"P", "B", "N", "R", "Q", "K", "p", "b", "n", "r", "q", "k"}
var PiecesUnicode = [12]string{"\u265F", "\u265D", "\u265E", "\u265C", "\u265B", "\u265A", "\u2659", "\u2657", "\u2658", "\u2656", "\u2655", "\u2654"}
var Sqaures = [2]string{"\u2B1B", "\u2B1C"}

func PieceSymbolToInt(piece string) uint8 {
	for i, p := range Pieces {
		if p == piece {
			return uint8(i) + 1
		}
	}
	return 0
}
func (b *Board) GetKing(color byte) (c Coord) {
	c.File = 8
	c.Rank = 8
	var king uint8 = 6
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

func (b *Board) GetPieces(color byte) (pieces []Coord) {
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

func GetPiece(b *Board, coord Coord) (piece uint8, color byte) {
	piece = b.AccessCoord(coord)
	if piece <= PieceOffset {
		color = WhiteToMove
	} else {
		color = BlackToMove
	}
	return
}

func (b *Board) GetLegalMoves(color byte) ([]string, []string) {
	m, c := b.GetMoves(color)
	return b.PruneIllegal(m, c)
}

func (b *Board) GetMoves(color byte) (moves, captures []string) {
	return b.getMoves(color, false)
}

func (b *Board) GetMovesNoCastling(color byte) (moves, captures []string) {
	return b.getMoves(color, true)
}

func (b *Board) getMoves(color byte, excludeCastling bool) (moves, captures []string) {
	pieces := b.GetPieces(color)
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

func (b *Board) placePiece(coord string, piece uint8) {
	target := AlgToCoord(coord)
	b.Coords[target.File][target.Rank] = piece
}
