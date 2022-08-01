package eval

import (
	"context"
	"fmt"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

func (e *EvalEngine) alphabetaWithTimeout(ctx context.Context, pv []board.Move, depth int, alpha, beta float32, side float32) (float32, []board.Move) {
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
			return side * e.EvalFunction(e, e.Board), make([]board.Move, 0)
		}

		var move board.Move
		var movesVar []board.Move
		var value float32

		value = negInf
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
			alpha = Max32(alpha, value)

			if alpha >= beta {
				break
			}

		}
		return value, append([]board.Move{move}, movesVar...)
	}
}

func (e *EvalEngine) IDSearch(ctx context.Context, depth int, alpha, beta float32) board.Move {
	var wg sync.WaitGroup
	var best board.Move
	var eval float32
	pv := make([]board.Move, depth)
	color := float32(-1)
	if e.Board.IsWhite {
		color = 1
	}
	done := false
	wg.Add(1)
	go func() {
		for d := 1; d <= depth; d++ {
			if done {
				wg.Done()
				return
			}

			tempEval, tempMove := e.alphabetaWithTimeout(ctx, pv, d, negInf, posInf, color)

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
				fmt.Printf("Depth: %d (%2.2f) Move: %v\n", d, eval, tempMove)
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
