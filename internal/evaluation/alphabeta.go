package eval

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

func (e *EvalEngine) negamax(ctx context.Context, line *[]board.Move, pvMoves []board.Move, depth int, alpha, beta int, side int) int {
	select {
	case <-ctx.Done():
		// Meaningless return. Should never trust the result after ctx is expired
		return 0
	default:
		if depth == 0 {
			return e.quiescence(ctx, alpha, beta, side)
		}

		e.Stats.nodes++

		// alphaTemp := alpha
		// if e.EnableTT {

		// 	if entry, ok := e.TTable[e.Board.Hash]; ok && entry.depth >= depth {
		// 		switch entry.ttType {
		// 		case TT_EXACT:
		// 			return entry.eval
		// 		case TT_LOWER:
		// 			alpha = Max(alpha, entry.eval)
		// 		case TT_UPPER:
		// 			beta = Min(beta, entry.eval)
		// 		}

		// 		if alpha >= beta {
		// 			*line = []board.Move{entry.move}
		// 			return entry.eval
		// 		}
		// 	}
		// }

		var pvMove board.Move
		if len(pvMoves) != 0 {
			pvMove = pvMoves[0]
			pvMoves = pvMoves[1:]
		}

		all := e.Board.PseudoMoveGen()
		legalMoves := 0
		e.Board.OrderMoves(pvMove, &all)

		var value int

		value = -math.MaxInt
		pv := []board.Move{}
		for i := 0; i < len(all); i++ {
			umove := e.Board.MakeMove(all[i])
			if e.Board.IsChecked(e.Board.Side ^ 1) {
				umove()
				continue
			}
			legalMoves++
			value = Max(value, -e.negamax(ctx, &pv, pvMoves, depth-1, -beta, -alpha, -side))
			umove()

			if value > alpha {
				alpha = value
				*line = []board.Move{all[i]}
				*line = append(*line, pv...)
			}

			if alpha >= beta {
				break
			}

			// if e.EnableTT {
			// 	tt := ttEntry{eval: value, depth: depth, move: all[i]}
			// 	if value <= alphaTemp {
			// 		tt.ttType = TT_UPPER
			// 	} else if value >= beta {
			// 		tt.ttType = TT_LOWER
			// 	} else {
			// 		tt.ttType = TT_EXACT
			// 	}
			// 	e.TTable[e.Board.Hash] = tt
			// }
		}

		if legalMoves == 0 {
			if e.Board.IsChecked(e.Board.Side) {
				return -math.MaxInt
			} else {
				return 0
			}
		}
		return value
	}
}

func (e *EvalEngine) quiescence(ctx context.Context, alpha, beta int, side int) int {
	select {
	case <-ctx.Done():
		// Meaningless return. Should never trust the result after ctx is expired
		return 0
	default:
		e.Stats.qNodes++
		eval := side * e.GetEvaluation(e.Board)

		if eval >= beta {
			return beta
		}

		if eval > alpha {
			alpha = eval
		}

		all := e.Board.PseudoCaptureGen()
		legalMoves := 0

		pvm := board.Move(0)
		e.Board.OrderMoves(pvm, &all)

		var value int

		value = -math.MaxInt
		for i := 0; i < len(all); i++ {
			umove := e.Board.MakeMove(all[i])
			if e.Board.IsChecked(e.Board.Side ^ 1) {
				umove()
				continue
			}
			legalMoves++
			value = Max(value, -e.quiescence(ctx, -beta, -alpha, -side))
			umove()
			alpha = Max(value, alpha)
			if alpha >= beta {
				break
			}
		}
		if legalMoves == 0 {
			return eval
		}
		return value
	}
}

// [100 325 325 500 975 10000]
// 2022/09/29 12:49:06 maxprocs: Leaving GOMAXPROCS=12: CPU quota undefined
// 2022/09/29 12:49:06 profile: cpu profiling enabled, /tmp/profile576563299/cpu.pprof
// Depth: 1 (0.20) Move: [b1c3] (403.8knps, total: 21.0 (1.0 20.0), QN: 95%, evals: 95%)
// Depth: 2 (0.09) Move: [b1c3 b8c6] (551.7knps, total: 112.0 (21.0 91.0), QN: 81%, evals: 81%)
// Depth: 3 (0.11) Move: [b1c3 b8c6 g1f3] (1.0Mnps, total: 571.0 (60.0 511.0), QN: 89%, evals: 89%)
// Depth: 4 (0.09) Move: [b1c3 b8c6 g1f3 g8f6] (671.7knps, total: 2.3k (524.0 1.8k), QN: 77%, evals: 77%)
// Depth: 5 (0.11) Move: [b1c3 b8c6 g1f3 g8f6 f3g1] (787.3knps, total: 16.9k (1.5k 15.4k), QN: 91%, evals: 91%)
// Depth: 6 (0.09) Move: [b1c3 b8c6 g1f3 g8f6 f3g1 f6g8] (431.0knps, total: 134.3k (13.8k 120.6k), QN: 89%, evals: 89%)
// Depth: 7 (0.11) Move: [b1c3 b8c6 g1f3 g8f6 f3g1 f6g8 g1f3] (740.3knps, total: 1.0M (52.2k 972.6k), QN: 94%, evals: 94%)
// Depth: 8 (0.09) Move: [b1c3 b8c6 g1f3 g8f6 f3g1 f6g8 g1f3 g8f6] (534.5knps, total: 10.8M (514.5k 10.3M), QN: 95%, evals: 95%)
// ^C2022/09/29 12:52:34 profile: caught interrupt, stopping profiles

// Iterative deepening search. Returns best move, ponder and ok if search succeeded.
func (e *EvalEngine) IDSearch(ctx context.Context, depth int, pv *[]board.Move, silent bool) (board.Move, board.Move, bool) {
	var wg sync.WaitGroup
	var best, ponder board.Move
	var eval int
	var line, bestLine []board.Move
	color := 1
	alpha, beta := -math.MaxInt, math.MaxInt
	if e.Board.Side != board.WHITE {
		color = -color
	}
	done, ok := false, true
	wg.Add(1)
	go func() {
		for d := 1; d <= depth; d++ {
			if done {
				wg.Done()
				return
			}

			if len(*pv) > len(line) {
				line = *pv
			}

			// e.TTable = make(map[uint64]ttEntry)
			e.Stats.Start()
			eval = e.negamax(ctx, &line, line, d, alpha, beta, color)

			select {
			case <-ctx.Done():
				// Do nothing as alpha-beta was canceled and results are unreliable
				done = true
				wg.Done()
				return
			default:
				// Debug purposes. No TTHit can happen under 4ply.
				// if depth < 4 && TTHits > 0 {
				// 	fmt.Println("TTable error - transposition in less than 4ply")
				// 	panic(1)
				// }

				if len(line) == 0 {
					done, ok = true, false
					break
				} else {
					best = line[0]
					if len(line) > 1 {
						ponder = line[1]
					}
					bestLine = line
				}
				if !silent {
					evalStr := ""
					switch eval {
					case math.MaxInt:
						evalStr = fmt.Sprintf("#%d", 1+len(line)/2)
					case -math.MaxInt:
						evalStr = fmt.Sprintf("#%d", 1+len(line)/2)
					default:
						evalStr = fmt.Sprintf("%2.2f", float32(color*eval)/100)
					}

					// fmt.Printf("Depth: %d (%s) Move: %v (TT hit: %d (Rate %2.2f%%) TT size: %d)\n", d, evalStr, line, TTHits, 100*float64(TTHits)/float64(len(e.TTable)), len(e.TTable))
					fmt.Printf("Depth: %d (%s) Move: %v (%s)\n", d, evalStr, line, e.Stats.String())
				}

				//found mate stop
				if eval == math.MaxInt || eval == -math.MaxInt {
					done = true
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()
	*pv = bestLine
	return best, ponder, ok
}
