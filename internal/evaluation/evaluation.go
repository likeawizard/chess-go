package eval

import (
	"context"
	"time"

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
type SearchFunction func(n *Node, depth ...int) float32
type EvalFunction func(*EvalEngine, *board.Board) float32

type EvalEngine struct {
	Evaluations   int64
	CachedEvals   int64
	EvalFunction  EvalFunction
	MoveTime      time.Duration
	Board         *board.Board
	DebugMode     bool
	SearchDepth   int
	TTable        map[uint64]float32
	MaxGoroutines chan struct{}
	Algorithm     string
}

func NewEvalEngine(b *board.Board, c *config.Config) (*EvalEngine, error) {
	return &EvalEngine{
		Board:         b,
		EvalFunction:  GetEvaluation,
		DebugMode:     c.Engine.Debug,
		SearchDepth:   c.Engine.MaxDepth,
		MaxGoroutines: make(chan struct{}, c.Engine.MaxGoRoutines),
		Algorithm:     c.Engine.Algorithm,
		TTable:        make(map[uint64]float32),
	}, nil
}

func (e *EvalEngine) GetMove(ctx context.Context) board.Move {
	e.Evaluations = 0
	start := time.Now()
	best := e.IDSearch(ctx, e.SearchDepth, negInf, posInf)
	e.MoveTime = time.Since(start)

	return best
}

func Max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func Min32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
