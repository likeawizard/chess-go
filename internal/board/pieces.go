package board

const PieceOffset = 6

var Pieces = [12]string{"P", "B", "N", "R", "Q", "K", "p", "b", "n", "r", "q", "k"}

func PieceSymbolToInt(piece string) uint8 {
	for i, p := range Pieces {
		if p == piece {
			return uint8(i) + 1
		}
	}
	return 0
}
func (b *Board) GetKing(isWhite bool) (c Square) {
	var king uint8 = 6
	if !isWhite {
		king += PieceOffset
	}

	for c = Square(0); c < 64; c++ {
		if b.Coords[c] == king {
			return
		}
	}
	return
}

func (b *Board) GetPieces(isWhite bool) (pieces []Square) {
	for i := Square(0); i < 64; i++ {
		piece := b.Coords[i]
		if piece == 0 {
			continue
		}

		if isWhite && piece < 7 || !isWhite && piece >= 7 {
			pieces = append(pieces, Square(i))
		}
	}
	return
}

func (b *Board) GetLegalMoves(isWhite bool) ([]Move, []Move) {
	m, c := b.GetMoves(isWhite)
	return b.PruneIllegal(m, c)
}

func (b *Board) GetMoves(isWhite bool) (moves, captures []Move) {
	return b.getMoves(isWhite, false)
}

func (b *Board) GetMovesNoCastling(isWhite bool) (moves, captures []Move) {
	return b.getMoves(isWhite, true)
}

func (b *Board) getMoves(isWhite bool, excludeCastling bool) (moves, captures []Move) {
	pieces := b.GetPieces(isWhite)
	for _, piece := range pieces {
		m, c := b.GetAvailableMovesRaw(piece, excludeCastling)
		moves = append(moves, m...)
		captures = append(captures, c...)
	}
	return
}
