package eval

import (
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/likeawizard/chess-go/internal/board"
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
	MoveToPlay string
	Evaluation float64
	Parent     *Node
	Children   []*Node
}
type SearchFunction func(n *Node, depth ...int) float64
type EvalFunction func(*EvalEngine, *board.Board) float64

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

func (n *Node) GetChildNodes() []*Node {
	fen := n.Position.ExportFEN()
	moves, captures := n.Position.GetMoves(n.Position.SideToMove)
	all := append(captures, moves...)
	childNodes := make([]*Node, len(all))

	for i := 0; i < len(all); i++ {
		childNodes[i] = &Node{
			Parent:     n,
			Children:   nil,
			Position:   &board.Board{},
			MoveToPlay: all[i],
		}
		childNodes[i].Position.ImportFEN(fen)
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

func NewEvalEngine(b *board.Board) (*EvalEngine, error) {
	debug, err := strconv.ParseBool(os.Getenv("EVALUATION_DEBUG"))
	if err != nil {
		return nil, err
	}
	depth, _ := strconv.Atoi(os.Getenv("EVALUATION_DEPTH"))
	if err != nil {
		return nil, err
	}
	max, _ := strconv.Atoi(os.Getenv("EVALUATION_MAX_GOROUTINES"))
	if err != nil {
		return nil, err
	}
	algo := os.Getenv("EVALUATION_ALGO")
	if err != nil {
		return nil, err
	}
	if algo != EVAL_ALPHABETA && algo != EVAL_MINMAX {
		algo = EVAL_ALPHABETA
	}
	return &EvalEngine{
		RootNode:      NewRootNode(b),
		EvalFunction:  SideDependantEval,
		DebugMode:     debug,
		SearchDepth:   depth,
		MaxGoroutines: make(chan struct{}, max),
		Algorithm:     algo,
	}, nil
}

func (e *EvalEngine) GetMove() {
	e.Evaluations = 0
	start := time.Now()
	switch e.Algorithm {
	case EVAL_MINMAX:
		e.minmax(e.RootNode, e.SearchDepth)
	case EVAL_ALPHABETA:
		e.alphabetaSerial(e.RootNode, e.SearchDepth, math.Inf(-1), math.Inf(1))
		// e.RootNode.alphabeta(e.SearchDepth, -15, 15)
	}
	e.MoveTime = time.Since(start)
}

func (n *Node) PickBestMove(side string) *Node {
	var bestMove *Node
	bestScore := math.Inf(-1)
	switch side {
	case board.WhiteToMove:
		for _, c := range n.Children {
			if c.Evaluation > bestScore {
				bestScore, bestMove = c.Evaluation, c
			}
		}
	case board.BlackToMove:
		bestScore = math.Inf(1)
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
	sort.Slice(moves, func(i, j int) bool {
		if n.Position.SideToMove == board.WhiteToMove {
			return moves[i].Evaluation > moves[j].Evaluation
		} else {
			return moves[i].Evaluation < moves[j].Evaluation
		}

	})
	return moves[:min(num, len(moves))]
}

func (n *Node) ConstructLine() []string {
	line := make([]string, 0)
	line = append(line, n.MoveToPlay)
	side := n.Position.SideToMove
	current := n
	for current.Children != nil {
		best := current.PickBestMove(side)
		line = append(line, best.MoveToPlay)
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

type CompFunc func(float64, float64) float64
type SelectiveCompFunc func(float64, float64, float64) (float64, float64)
type CompFuncBool func(float64, float64, float64) bool

func gte(x, a, b float64) bool {
	return x >= b
}

func lte(x, a, b float64) bool {
	return x <= a
}

func minB(x, a, b float64) (float64, float64) {
	return a, math.Min(x, b)
}

func maxA(x, a, b float64) (float64, float64) {
	return math.Max(x, a), b
}
