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
func (b *Board) GetKing(color byte) (c Square) {
	var king uint8 = 6
	if color == BlackToMove {
		king += PieceOffset
	}

	for c = Square(0); c < 64; c++ {
		if b.Coords[c] == king {
			return
		}
	}
	return
}

func (b *Board) GetPieces(color byte) (pieces []Square) {
	for i := Square(0); i < 64; i++ {
		piece := b.Coords[i]
		if piece == 0 {
			continue
		}

		if color == WhiteToMove && piece < 7 || color == BlackToMove && piece >= 7 {
			pieces = append(pieces, Square(i))
		}
	}
	return
}

func (b *Board) GetLegalMoves(color byte) ([]Move, []Move) {
	m, c := b.GetMoves(color)
	return b.PruneIllegal(m, c)
}

func (b *Board) GetMoves(color byte) (moves, captures []Move) {
	return b.getMoves(color, false)
}

func (b *Board) GetMovesNoCastling(color byte) (moves, captures []Move) {
	return b.getMoves(color, true)
}

func (b *Board) getMoves(color byte, excludeCastling bool) (moves, captures []Move) {
	pieces := b.GetPieces(color)
	for _, piece := range pieces {
		m, c := b.GetAvailableMovesRaw(piece, excludeCastling)
		moves = append(moves, m...)
		captures = append(captures, c...)
	}
	return
}
