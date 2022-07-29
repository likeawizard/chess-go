package eval

import (
	"context"
	"fmt"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

var ac, bc int

func (e *EvalEngine) alphabetaWithTimeout(ctx context.Context, depth int, alpha, beta float32) (float32, board.Move) {
	select {
	case <-ctx.Done():
		// Meaningless return. Should never trust the result after ctx is expired
		return 0, 0
	default:
		if depth == 0 {
			return e.EvalFunction(e, e.Board), 0
		}

		m, c := e.Board.GetLegalMoves()
		all := append(c, m...)
		var move board.Move
		var value, temp float32
		if e.Board.IsWhite {
			value = negInf
			for i := 0; i < len(all); i++ {
				umove := e.Board.MoveLongAlg(all[i])
				temp, _ = e.alphabetaWithTimeout(ctx, depth-1, alpha, beta)
				if temp > value {
					move = all[i]
				}
				value = Max32(value, temp)
				umove()

				if value >= beta {
					bc++
					break
				}
				alpha = Max32(alpha, value)
			}
			return value, move
		} else {
			value = posInf
			for i := 0; i < len(all); i++ {
				umove := e.Board.MoveLongAlg(all[i])
				temp, _ = e.alphabetaWithTimeout(ctx, depth-1, alpha, beta)
				if temp < value {
					move = all[i]
				}
				value = Min32(value, temp)
				umove()

				if value <= alpha {
					ac++
					break
				}
				beta = Min32(beta, value)
			}
			return value, move
		}
	}
}

func (e *EvalEngine) IDSearch(ctx context.Context, depth int, alpha, beta float32) board.Move {
	var wg sync.WaitGroup
	var best board.Move
	var eval float32
	done := false
	wg.Add(1)
	go func() {
		for d := 1; d <= depth; d++ {
			if done {
				wg.Done()
				return

			}

			ac, bc = 0, 0
			tempEval, tempMove := e.alphabetaWithTimeout(ctx, d, negInf, posInf)

			if tempMove == 0 {
				break
			}

			select {
			case <-ctx.Done():
				// Do nothing as alpha-beta was canceled and results are unreliable
				done = true
				wg.Done()
				return
			default:
				eval, best = tempEval, tempMove
				fmt.Printf("Depth: %d Move: %v (%2.2f)(ac: %d bc: %d)\n", d, best, eval, ac, bc)
				//found mate stop
				if tempEval == negInf || tempEval == posInf {
					done = true
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()
	return best
}
