package eval

import (
	"context"
	"fmt"
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

type EvalFunction func(*EvalEngine, *board.Board) float32

type EvalEngine struct {
	Evaluations   int64
	CachedEvals   int64
	EvalFunction  EvalFunction
	MoveTime      time.Duration
	RootNode      *Node
	DebugMode     bool
	SearchDepth   int
	MaxGoroutines chan struct{}
	Algorithm     string
}

func NewEvalEngine(b *board.Board, c *config.Config) (*EvalEngine, error) {
	return &EvalEngine{
		RootNode:      NewRootNode(b),
		EvalFunction:  GetEvaluation,
		DebugMode:     c.Engine.Debug,
		SearchDepth:   c.Engine.MaxDepth,
		MaxGoroutines: make(chan struct{}, c.Engine.MaxGoRoutines),
		Algorithm:     c.Engine.Algorithm,
	}, nil
}

func (e *EvalEngine) BuildGameTree(depth int) {
	root := e.RootNode
	root.BuildGameTree(depth)
}

func (e *EvalEngine) EvaluateLeafNodes(ctx context.Context) {
	go func(ctx context.Context) {
		e.RootNode.EvaluateLeafNodes(ctx, e)
	}(ctx)
}

func (e *EvalEngine) PonderOnMove(ctx context.Context) {
	e.BuildGameTree(e.SearchDepth + 1)
	// e.EvaluateLeafNodes(ctx)
	// fmt.Println("built search tree")
}

func (e *EvalEngine) ResetRootWithMove(move string) error {
	for _, child := range e.RootNode.Children {
		if child.MoveToPlay == move {
			e.RootNode = child
			return nil
		}
	}
	return fmt.Errorf("move not found among children: %s", move)
}

func (e *EvalEngine) GetMove() {
	fmt.Println(e.RootNode.Position.ExportFEN())
	e.Evaluations = 0
	start := time.Now()
	switch e.Algorithm {
	case EVAL_MINMAX:
		e.minmaxSerial(e.RootNode, e.SearchDepth, e.RootNode.Position.SideToMove == board.WhiteToMove)
	case EVAL_ALPHABETA:
		e.alphabetaSerial(e.RootNode, e.SearchDepth, negInf, posInf, e.RootNode.Position.SideToMove == board.WhiteToMove)
	}
	// for _, c := range e.RootNode.Children {
	// 	fmt.Println(c.MoveToPlay)
	// }

	e.MoveTime = time.Since(start)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (e *EvalEngine) PlayMove(move *Node) {
	e.RootNode = move
	e.RootNode.Parent = nil
}

type CompFunc func(float32, float32) float32
type SelectiveCompFunc func(float32, float32, float32) (float32, float32)
type CompFuncBool func(float32, float32, float32) bool

func gte(x, a, b float32) bool {
	return x >= b
}

func lte(x, a, b float32) bool {
	return x <= a
}

func minB(x, a, b float32) (float32, float32) {
	return a, Min32(x, b)
}

func maxA(x, a, b float32) (float32, float32) {
	return Max32(x, a), b
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
