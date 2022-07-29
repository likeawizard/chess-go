package eval

import (
	"context"
	"fmt"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

var ac, bc int

func (e *EvalEngine) alphabetaWithTimeout(ctx context.Context, pv []board.Move, depth int, alpha, beta float32) (float32, []board.Move) {
	select {
	case <-ctx.Done():
		// Meaningless return. Should never trust the result after ctx is expired
		return 0, nil
	default:
		m, c := e.Board.GetLegalMoves()
		pvm := board.Move(0)
		if len(pv) > 0 {
			pvm = pv[0]
			pv = pv[1:]
		}
		all := e.Board.OrderMoves(pvm, m, c)
		if depth == 0 || len(all) == 0 {
			return e.EvalFunction(e, e.Board), make([]board.Move, 0)
		}

		var move board.Move
		var movesVar []board.Move
		var value float32
		if e.Board.IsWhite {
			value = negInf
			move = all[0]
			for i := 0; i < len(all); i++ {
				umove := e.Board.MoveLongAlg(all[i])
				temp, tempMoves := e.alphabetaWithTimeout(ctx, pv, depth-1, alpha, beta)
				if temp > value {
					value = temp
					move = all[i]
					movesVar = tempMoves
				}
				umove()

				if value >= beta {
					bc++
					break
				}
				alpha = Max32(alpha, value)
			}
			return value, append([]board.Move{move}, movesVar...)
		} else {
			value = posInf
			move = all[0]
			for i := 0; i < len(all); i++ {
				umove := e.Board.MoveLongAlg(all[i])
				temp, tempMoves := e.alphabetaWithTimeout(ctx, pv, depth-1, alpha, beta)
				if temp < value {
					value = temp
					move = all[i]
					movesVar = tempMoves
				}
				umove()

				if value <= alpha {
					ac++
					break
				}
				beta = Min32(beta, value)
			}
			return value, append([]board.Move{move}, movesVar...)
		}
	}
}

func (e *EvalEngine) IDSearch(ctx context.Context, depth int, alpha, beta float32) board.Move {
	var wg sync.WaitGroup
	var best board.Move
	var eval float32
	pv := make([]board.Move, depth)
	done := false
	wg.Add(1)
	go func() {
		for d := 1; d <= depth; d++ {
			if done {
				wg.Done()
				return

			}

			ac, bc = 0, 0
			tempEval, tempMove := e.alphabetaWithTimeout(ctx, pv, d, negInf, posInf)

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
				eval, best = tempEval, tempMove[0]
				pv = tempMove
				fmt.Printf("Depth: %d (%2.2f) Move: %v (ac: %d bc: %d)\n", d, eval, tempMove, ac, bc)
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
