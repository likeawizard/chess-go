package eval

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

var TTHits int

func (e *EvalEngine) negamax(ctx context.Context, line *[]board.Move, pvMoves []board.Move, depth int, alpha, beta int, side int) int {
	select {
	case <-ctx.Done():
		// Meaningless return. Should never trust the result after ctx is expired
		return 0
	default:
		if depth == 0 {
			return e.quiescence(ctx, alpha, beta, side)
		}

		alphaTemp := alpha
		if e.EnableTT {

			if entry, ok := e.TTable[e.Board.Hash]; ok && entry.depth >= depth {
				TTHits++
				switch entry.ttType {
				case TT_EXACT:
					return entry.eval
				case TT_LOWER:
					alpha = Max(alpha, entry.eval)
				case TT_UPPER:
					beta = Min(beta, entry.eval)
				}

				if alpha >= beta {
					*line = []board.Move{entry.move}
					return entry.eval
				}
			}
		}

		var pvMove board.Move
		if len(pvMoves) != 0 {
			pvMove = pvMoves[0]
			pvMoves = pvMoves[1:]
		}

		all := e.Board.MoveGen()
		e.Board.OrderMoves(pvMove, &all)

		if len(all) == 0 {
			return side * e.GetEvaluation(e.Board)
		}

		var value int

		value = -math.MaxInt
		pv := []board.Move{}
		for i := 0; i < len(all); i++ {
			umove := e.Board.MakeMove(all[i])
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

			if e.EnableTT {
				tt := ttEntry{eval: value, depth: depth, move: all[i]}
				if value <= alphaTemp {
					tt.ttType = TT_UPPER
				} else if value >= beta {
					tt.ttType = TT_LOWER
				} else {
					tt.ttType = TT_EXACT
				}
				e.TTable[e.Board.Hash] = tt
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
		eval := side * e.GetEvaluation(e.Board)

		if eval >= beta {
			return beta
		}

		if eval > alpha {
			alpha = eval
		}

		all := e.Board.MoveGenCaptures()
		pvm := board.Move(0)
		e.Board.OrderMoves(pvm, &all)

		if len(all) == 0 {
			return eval
		}

		var value int

		value = -math.MaxInt
		for i := 0; i < len(all); i++ {
			umove := e.Board.MakeMove(all[i])
			value = Max(value, -e.quiescence(ctx, -beta, -alpha, -side))
			umove()
			alpha = Max(value, alpha)
			if alpha >= beta {
				break
			}
		}
		return value
	}
}

// Iterative deepening search. Returns best move, ponder and ok if search succeeded.
func (e *EvalEngine) IDSearch(ctx context.Context, depth int, pv *[]board.Move, silent bool) (board.Move, board.Move, bool) {
	var wg sync.WaitGroup
	var best, ponder board.Move
	var eval int
	var line, bestLine []board.Move
	color := 1
	alpha, beta := -math.MaxInt, math.MaxInt
	if !e.Board.IsWhite {
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

			e.TTable = make(map[uint64]ttEntry)
			TTHits = 0
			eval = e.negamax(ctx, &line, line, d, alpha, beta, color)

			select {
			case <-ctx.Done():
				// Do nothing as alpha-beta was canceled and results are unreliable
				done = true
				wg.Done()
				return
			default:
				// Debug purposes. No TTHit can happen under 4ply.
				if depth < 4 && TTHits > 0 {
					fmt.Println("TTable error - transposition in less than 4ply")
					panic(1)
				}

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

					fmt.Printf("Depth: %d (%s) Move: %v (TT hit: %d (Rate %2.2f%%) TT size: %d)\n", d, evalStr, line, TTHits, 100*float64(TTHits)/float64(len(e.TTable)), len(e.TTable))
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
