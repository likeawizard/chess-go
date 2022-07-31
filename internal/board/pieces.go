package board

import (
	"sort"
)

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
	var king uint8 = K
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

func (b *Board) GetLegalMoves() (moves, captures []Move) {
	pins := b.GetPins(b.IsWhite)
	_ = pins

	var checks []Move
	check := Move(0)
	if inCheck := b.IsInCheck(b.IsWhite); inCheck {
		checks = b.GetChecks(b.IsWhite)
		check = checks[0]
	}

	// If double-check only consider king moves
	if len(checks) == 2 {
		return b.GetMovesForPiece(b.GetKing(b.IsWhite), 0, 0)
	} else {
		pieces := b.GetPieces(b.IsWhite)
		for _, piece := range pieces {
			pin := getPin(piece, pins)
			m, c := b.GetMovesForPiece(piece, pin, check)
			moves = append(moves, m...)
			captures = append(captures, c...)
		}
		return
	}
}

func getPin(sq Square, pins []Move) Move {
	for _, pin := range pins {
		if sq == pin.From() {
			return pin
		}
	}
	return 0
}

func (b *Board) OrderMoves(pv Move, moves, captures []Move) []Move {
	all := append(captures, moves...)

	sort.Slice(all, func(i int, j int) bool {
		return all[i] == pv || b.getMoveValue(all[i]) > b.getMoveValue(all[j])
	})

	return all
}

var PieceWeights = [13]float32{0, 1, 3.2, 2.9, 5, 9, 0, -1, -3.2, -2.9, -5, -9, 0}

func (b *Board) getMoveValue(capture Move) float32 {
	dir := float32(-1)
	if !b.IsWhite {
		dir *= -1
	}
	from, to := capture.FromTo()
	us, them := PieceWeights[b.Coords[from]], PieceWeights[b.Coords[to]]
	if them == 0 {
		return 0
	} else {
		return dir * (0.5*us + them)
	}
}
