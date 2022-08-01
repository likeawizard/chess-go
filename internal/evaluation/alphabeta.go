package eval

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

func (e *EvalEngine) alphabetaWithTimeout(ctx context.Context, pv []board.Move, depth int, alpha, beta int, side int) (int, []board.Move) {
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
				eval, best = tempEval, tempMove[0]
				pv = tempMove
				fmt.Printf("Depth: %d (%2.2f) Move: %v\n", d, float32(eval)/100, tempMove)
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
