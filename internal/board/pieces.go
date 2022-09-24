package board

import (
	"sort"
)

var Pieces = [6]string{"P", "B", "N", "R", "Q", "K"}

func (b *Board) OrderMoves(pv Move, moves *[]Move) {
	sort.Slice(*moves, func(i int, j int) bool {
		return (*moves)[i] == pv || b.getMoveValue((*moves)[i]) > b.getMoveValue((*moves)[j])
	})
}

// Estimate the potential strength of the move for move ordering
func (b *Board) getMoveValue(move Move) (value float32) {

	dir := float32(-1)
	if b.Side != WHITE {
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
