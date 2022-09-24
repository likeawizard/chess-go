package eval

import (
	"context"

	"github.com/likeawizard/chess-go/internal/board"
	"github.com/likeawizard/chess-go/internal/config"
)

type EvalEngine struct {
	Stats       EvalStats
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
	var best, ponder board.Move
	var ok bool
	all := e.Board.MoveGen()
	if len(all) == 1 {
		best = all[0]
	} else {
		best, ponder, ok = e.IDSearch(ctx, e.SearchDepth, pv, silent)
		if !ok {
			best = all[0]
		}
	}

	return best, ponder
}

// TODO: try branchless optimization
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
