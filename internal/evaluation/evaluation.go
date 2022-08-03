package eval

import (
	"context"

	"github.com/likeawizard/chess-go/internal/board"
	"github.com/likeawizard/chess-go/internal/config"
)

var DEBUG = false
var MAX_DEPTH int

var Evaluations int64
var CachedEvals int64

const (
	EVAL_MINMAX    string = "minmax"
	EVAL_ALPHABETA string = "alphabeta"
)

type Node struct {
	Position   *board.Board
	MoveToPlay board.Move
	Evaluation float32
	Parent     *Node
	Children   []*Node
}
type SearchFunction func(n *Node, depth ...int) int
type EvalFunction func(*EvalEngine, *board.Board) int

type EvalEngine struct {
	Evaluations  int64
	EvalFunction EvalFunction
	Board        *board.Board
	DebugMode    bool
	SearchDepth  int
	EnableTT     bool
	TTable       map[uint64]ttEntry
}

func NewEvalEngine(b *board.Board, c *config.Config) (*EvalEngine, error) {
	return &EvalEngine{
		Board:        b,
		EvalFunction: GetEvaluation,
		SearchDepth:  c.Engine.MaxDepth,
		EnableTT:     c.Engine.EnableTT,
		TTable:       make(map[uint64]ttEntry),
	}, nil
}

func (e *EvalEngine) GetMove(ctx context.Context) board.Move {
	e.Evaluations = 0
	var best board.Move
	m, c := e.Board.GetLegalMoves()
	all := append(m, c...)
	if len(all) == 1 {
		best = all[0]
	} else {
		best = e.IDSearch(ctx, e.SearchDepth)
	}

	return best
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
