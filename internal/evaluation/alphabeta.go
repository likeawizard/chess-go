package eval

import (
	"context"
	"sort"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

func (e *EvalEngine) alphabetaSerial(n *Node, depth int, alpha, beta float32, isWhite bool) float32 {
	if depth == 0 {
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	}

	if n.Children == nil {
		n.Children = n.GetChildNodes()
	}

	var value float32
	if isWhite {
		value = negInf
		for i := 0; i < len(n.Children); i++ {

			value = Max32(value, e.alphabetaSerial(n.Children[i], depth-1, alpha, beta, false))

			if value >= beta {
				break
			}
			alpha = Max32(alpha, value)
		}
		if n != nil {
			n.Evaluation = value
		}

		return value
	} else {
		value = posInf
		for i := 0; i < len(n.Children); i++ {
			value = Min32(value, e.alphabetaSerial(n.Children[i], depth-1, alpha, beta, true))

			if value <= alpha {
				break
			}
			beta = Min32(beta, value)
		}
		if n != nil {
			n.Evaluation = value
		}
		return value
	}
}

func (e *EvalEngine) alphabetaSerialWithTimeout(ctx context.Context, n *Node, depth int, alpha, beta float32, isWhite bool) float32 {
	select {
	case <-ctx.Done():
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	default:
		if val, ok := e.TTable[n.Position.Hash]; ok {
			// TThit++
			// storeVariations(n, n.Position.Hash)
			return val
		}
		if depth == 0 {
			n.Evaluation = e.EvalFunction(e, n.Position)
			return n.Evaluation
		}

		if n.Children == nil {
			n.Children = n.GetChildNodes()
		}

		var value float32
		if isWhite {
			value = negInf
			for i := 0; i < len(n.Children); i++ {

				value = Max32(value, e.alphabetaSerialWithTimeout(ctx, n.Children[i], depth-1, alpha, beta, false))

				if value >= beta {
					break
				}
				alpha = Max32(alpha, value)
			}

			e.TTable[n.Position.Hash] = value
			// storeVariations(n, n.Position.Hash)
			n.Evaluation = value
			return value
		} else {
			value = posInf
			for i := 0; i < len(n.Children); i++ {
				value = Min32(value, e.alphabetaSerialWithTimeout(ctx, n.Children[i], depth-1, alpha, beta, true))

				if value <= alpha {
					break
				}
				beta = Min32(beta, value)
			}

			e.TTable[n.Position.Hash] = value
			// storeVariations(n, n.Position.Hash)
			n.Evaluation = value
			return value
		}
	}
}

// func storeVariations(n *Node, hash uint64) {
// 	variation := ""
// 	for n.Parent != nil {
// 		variation = n.MoveToPlay + " " + variation
// 		n = n.Parent
// 	}
// 	TTDebug[hash] = append(TTDebug[hash], variation)
// }

// var TThit uint64

// var TTDebug map[uint64][]string

// func printTrans() {
// 	for hash, vari := range TTDebug {
// 		if len(vari) > 1 {
// 			fmt.Printf("%d: %v\n", hash, vari)
// 		}
// 	}
// }

func (e *EvalEngine) alphaBetaWithOrdering(ctx context.Context, n *Node, depth int, alpha, beta float32, isWhite bool) *Node {
	var best *Node
	var wg sync.WaitGroup
	// start := time.Now()
	wg.Add(1)
	go func() {
		for d := 1; d <= depth; d++ {
			e.TTable = make(map[uint64]float32)
			// TTDebug = make(map[uint64][]string)
			// TThit = 0
			// e.Evaluations = 0
			// currTime := time.Now()
			e.alphabetaSerialWithTimeout(ctx, n, d, alpha, beta, isWhite)
			select {
			case <-ctx.Done():
				// Do nothing as alpha-beta was canceled and results are unreliable
				// fmt.Printf("Timeout. Canceling search at depth: %d (time spent: %v)\n", d, time.Since(start))
				wg.Done()
				return
			default:
				best = n.PickBestMove(n.Position.SideToMove)
				nodes := e.orderTree(n, 2)
				_ = nodes
				// bFactor := math.Pow(float64(nodes), 1.0/float64(d))
				// fmt.Printf("Best move: %s. Found at depth: %d (branching factor: %v time spent: total = %v at depth = %v estimate on next step: %v)\n", best.MoveToPlay, d, bFactor, time.Since(start), time.Since(currTime), time.Duration(time.Since(start)*time.Duration(bFactor)))
				// fmt.Printf("Best move: %s. Found at depth: %d (time spent: total = %v at depth = %v) (Nodes: %d TT Hits:%d TT Size: %d)\n", best.MoveToPlay, d, time.Since(start), time.Since(currTime), e.Evaluations, TThit, len(e.TTable))
				// printTrans()
			}
		}
		wg.Done()
	}()

	wg.Wait()
	return best

}

func (e *EvalEngine) orderTree(n *Node, depth int) int {
	if depth == 0 {
		return 0
	}
	nodes := len(n.Children)
	sort.Slice(n.Children, func(i, j int) bool {
		if n.Position.SideToMove == board.WhiteToMove {
			return n.Children[i].Evaluation > n.Children[j].Evaluation
		} else {
			return n.Children[i].Evaluation < n.Children[j].Evaluation
		}
	})
	for i := 0; i < len(n.Children); i++ {
		nodes += e.orderTree(n.Children[i], depth-1)
	}
	return nodes
}
