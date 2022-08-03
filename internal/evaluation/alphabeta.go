package eval

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

var TTHits int

func (e *EvalEngine) alphabetaWithTimeout(ctx context.Context, pv []board.Move, depth int, alpha, beta int, side int) (int, []board.Move) {
	select {
	case <-ctx.Done():
		// Meaningless return. Should never trust the result after ctx is expired
		return 0, nil
	default:
		alphaTemp := alpha

		if entry, ok := e.TTable[e.Board.Hash]; ok && entry.depth >= depth {
			TTHits++
			switch entry.ttType {
			case TT_EXACT:
				return entry.eval, make([]board.Move, 0)
			case TT_LOWER:
				alpha = Max(alpha, entry.eval)
			case TT_UPPER:
				beta = Min(beta, entry.eval)
			}

			if alpha >= beta {
				return entry.eval, make([]board.Move, 0)
			}
		}
		m, c := e.Board.GetLegalMoves()
		pvm := board.Move(0)
		if len(pv) > 0 {
			pvm = pv[0]
			pv = pv[1:]
		}
		all := e.Board.OrderMoves(pvm, m, c)

		if depth == 0 || len(all) == 0 {
			return side * e.EvalFunction(e, e.Board), make([]board.Move, 0)
		}

		var move board.Move
		var movesVar []board.Move
		var value int

		value = -math.MaxInt
		move = all[0]
		for i := 0; i < len(all); i++ {
			umove := e.Board.MoveLongAlg(all[i])
			temp, tempMoves := e.alphabetaWithTimeout(ctx, pv, depth-1, -beta, -alpha, -side)
			temp = -temp
			if temp > value {
				value = temp
				move = all[i]
				movesVar = tempMoves
			}
			umove()
			alpha = Max(alpha, value)

			if alpha >= beta {
				break
			}

			tt := ttEntry{eval: value, depth: depth}
			if value <= alphaTemp {
				tt.ttType = TT_UPPER
			} else if value >= beta {
				tt.ttType = TT_LOWER
			} else {
				tt.ttType = TT_EXACT
			}
			e.TTable[e.Board.Hash] = tt
		}
		return value, append([]board.Move{move}, movesVar...)
	}
}

func (e *EvalEngine) IDSearch(ctx context.Context, depth int) board.Move {
	var wg sync.WaitGroup
	var best board.Move
	var eval int
	pv := make([]board.Move, depth)
	color := 1
	alpha, beta := -math.MaxInt, math.MaxInt
	if !e.Board.IsWhite {
		color = -color
	}
	done := false
	wg.Add(1)
	go func() {
		for d := 1; d <= depth; d++ {
			if done {
				wg.Done()
				return
			}
			e.TTable = make(map[uint64]ttEntry)
			TTHits = 0
			tempEval, tempMove := e.alphabetaWithTimeout(ctx, pv, d, alpha, beta, color)

			if len(tempMove) == 0 {
				break
			}

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
				eval, best = tempEval, tempMove[0]
				pv = tempMove
				fmt.Printf("Depth: %d (%2.2f) Move: %v (TT hit: %d (Rate %2.2f%%) TT size: %d)\n", d, float32(color*eval)/100, tempMove, TTHits, 100*float64(TTHits)/float64(len(e.TTable)), len(e.TTable))
				//found mate stop
				if tempEval == math.MaxInt || tempEval == -math.MaxInt {
					done = true
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()
	return best
}
