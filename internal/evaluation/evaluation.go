package eval

import (
	"context"
	"fmt"

	"github.com/likeawizard/chess-go/internal/board"
	"github.com/likeawizard/chess-go/internal/config"
)

type EvalEngine struct {
	Evaluations int64
	Board       *board.Board
	SearchDepth int
	EnableTT    bool
	TTable      map[uint64]ttEntry
}

func NewEvalEngine(b *board.Board, c *config.Config) (*EvalEngine, error) {
	return &EvalEngine{
		Board:       b,
		SearchDepth: c.Engine.MaxDepth,
		EnableTT:    c.Engine.EnableTT,
		TTable:      make(map[uint64]ttEntry),
	}, nil
}

// Returns the best move and best opponent response - ponder
func (e *EvalEngine) GetMove(ctx context.Context, pv *[]board.Move, silent bool) (board.Move, board.Move) {
	e.Evaluations = 0
	var best, ponder board.Move
	var ok bool
	all := e.Board.GetLegalMoves()
	if len(all) == 1 {
		best = all[0]
	} else {
		best, ponder, ok = e.IDSearch(ctx, e.SearchDepth, pv, silent)
		if !ok {
			best = all[0]
		}
	}
	fmt.Println(*pv)
	return best, ponder
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
