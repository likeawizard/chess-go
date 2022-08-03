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
		all := b.GetLegalMoves()
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
	all := b.GetLegalMoves()

	nodesSearched := 0
	for _, move := range all {
		unmove := b.MoveLongAlg(move)
		nodes := traverse(b, depth-1)
		nodesSearched += nodes
		fmt.Printf("%s: %d\n", move, nodes)
		unmove()
	}
	fmt.Println("\nNodes searched: ", nodesSearched)
}
