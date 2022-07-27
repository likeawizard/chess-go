package eval

import (
	"context"
	"fmt"
	"sort"
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
	RootNode      *Node
	DebugMode     bool
	SearchDepth   int
	TTable        map[uint64]float32
	MaxGoroutines chan struct{}
	Algorithm     string
}

func (n *Node) GetChildNodes() []*Node {
	moves, captures := n.Position.GetLegalMoves(n.Position.SideToMove)
	all := append(captures, moves...)
	childNodes := make([]*Node, len(all))

	for i := 0; i < len(all); i++ {
		childNodes[i] = &Node{
			Parent:     n,
			Children:   nil,
			Position:   &board.Board{},
			MoveToPlay: all[i],
		}
		childNodes[i].Position = n.Position.SimpleCopy()
		childNodes[i].Position.MoveLongAlg(all[i])
	}

	return childNodes
}

func NewRootNode(b *board.Board) *Node {
	node := &Node{
		Position: b,
		Parent:   nil,
	}

	node.Children = node.GetChildNodes()

	return node
}

func NewEvalEngine(b *board.Board, c *config.Config) (*EvalEngine, error) {
	return &EvalEngine{
		RootNode:      NewRootNode(b),
		EvalFunction:  GetEvaluation,
		DebugMode:     c.Engine.Debug,
		SearchDepth:   c.Engine.MaxDepth,
		MaxGoroutines: make(chan struct{}, c.Engine.MaxGoRoutines),
		Algorithm:     c.Engine.Algorithm,
		TTable:        make(map[uint64]float32),
	}, nil
}

func (e *EvalEngine) GetMove(ctx context.Context) *Node {
	var best *Node
	e.Evaluations = 0
	start := time.Now()
	switch e.Algorithm {
	case EVAL_MINMAX:
		e.minmaxSerial(e.RootNode, e.SearchDepth, e.RootNode.Position.SideToMove == board.WhiteToMove)
		best = e.RootNode.PickBestMove(e.RootNode.Position.SideToMove)
	case EVAL_ALPHABETA:
		best = e.alphaBetaWithOrdering(ctx, e.RootNode, e.SearchDepth, negInf, posInf, e.RootNode.Position.SideToMove == board.WhiteToMove)
	}
	e.MoveTime = time.Since(start)

	return best
}

func (n *Node) PickBestMove(side byte) *Node {
	if n.Children == nil || len(n.Children) == 0 {
		return nil
	}
	var bestMove *Node = n.Children[0]
	bestScore := negInf
	switch side {
	case board.WhiteToMove:
		for _, c := range n.Children {
			if c.Evaluation > bestScore {
				bestScore, bestMove = c.Evaluation, c
			}
		}
	case board.BlackToMove:
		bestScore = posInf
		for _, c := range n.Children {
			if c.Evaluation < bestScore {
				bestScore, bestMove = c.Evaluation, c
			}
		}
	}

	return bestMove
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (n *Node) PickBestMoves(num int) []*Node {
	moves := n.Children
	num = min(num, len(moves))
	sort.Slice(moves, func(i, j int) bool {
		if n.Position.SideToMove == board.WhiteToMove {
			return moves[i].Evaluation > moves[j].Evaluation
		} else {
			return moves[i].Evaluation < moves[j].Evaluation
		}

	})
	return moves[:num]
}

func (n *Node) ConstructLine() []string {
	line := make([]string, 0)
	line = append(line, n.MoveToPlay.String())
	side := n.Position.SideToMove
	current := n
	for current.Children != nil {
		best := current.PickBestMove(side)
		if best == nil {
			break
		}
		line = append(line, best.MoveToPlay.String())
		switch side {
		case board.WhiteToMove:
			side = board.BlackToMove
		case board.BlackToMove:
			side = board.WhiteToMove
		}
		current = best
	}

	return line
}

func (e *EvalEngine) PlayMove(move *Node) {
	e.RootNode = move
	e.RootNode.Parent = nil
}

func (e *EvalEngine) ResetRootWithMove(move string) error {
	for _, child := range e.RootNode.Children {
		if child.MoveToPlay.String() == move {
			e.RootNode = child
			return nil
		}
	}
	return fmt.Errorf("move not found among children: %s", move)
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
