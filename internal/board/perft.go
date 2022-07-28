package board

import (
	"fmt"
	"time"
)

func Perft(fen string, depth int) (int, time.Duration) {
	b := &Board{}
	b.ImportFEN(fen)
	start := time.Now()
	leafs := traverse(b, depth)
	return leafs, time.Since(start)
}

func traverse(b *Board, depth int) int {
	num := 0
	if depth == 0 {
		return 1
	} else {
		m, c := b.GetLegalMoves(b.IsWhite)
		all := append(m, c...)
		for i := 0; i < len(all); i++ {
			umove := b.MoveLongAlg(all[i])
			num += traverse(b, depth-1)
			umove()
		}
		return num
	}
}

func PerftDebug(fen string, depth int) {
	b := &Board{}
	b.ImportFEN(fen)
	m, c := b.GetLegalMoves(b.IsWhite)
	all := append(m, c...)
	for _, move := range all {
		unmove := b.MoveLongAlg(move)
		fmt.Printf("%s: %d\n", move, traverse(b, depth-1))
		unmove()
	}
}
