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

func (b *Board) OrderMoves(pv Move, moves *[]Move) {
	sort.Slice(*moves, func(i int, j int) bool {
		return (*moves)[i] == pv || b.getMoveValue((*moves)[i]) > b.getMoveValue((*moves)[j])
	})
}

var PieceWeights = [13]float32{0, 1, 3.2, 2.9, 5, 9, 0, -1, -3.2, -2.9, -5, -9, 0}

// Estimate the potential strength of the move for move ordering
func (b *Board) getMoveValue(move Move) (value float32) {

	dir := float32(-1)
	if !b.IsWhite {
		dir *= -1
	}

	// TODO: implement SEE or MVV-LVA ordering
	// Calculate the relative value of exchange
	// from, to := move.FromTo()
	// us, them := PieceWeights[b.Coords[from]], PieceWeights[b.Coords[to]]
	// if them == 0 {
	// 	value += 0
	// } else {
	// 	value += dir * (0.5*us + them)
	// }

	// Prioritize promotions
	if move.Promotion() != 0 {
		value += 3
	}

	return
}
