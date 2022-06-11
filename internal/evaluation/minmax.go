package eval

import (
	"math"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

func (e *EvalEngine) minmax(n *Node, depth int) float64 {
	if depth == 0 {
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	}

	if n.Children == nil {
		n.Children = n.GetChildNodes()
	}

	val := math.Inf(-1)
	var comp CompFunc = math.Max

	if n.Position.SideToMove == board.BlackToMove {
		val = math.Inf(1)
		comp = math.Min
	}

	var wg sync.WaitGroup
	temp := make([]float64, len(n.Children))

	for i := 0; i < len(n.Children); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			temp[i] = e.minmaxSerial(n.Children[i], depth-1)
		}(i)
	}

	wg.Wait()

	for i := 0; i < len(temp); i++ {
		val = comp(val, temp[i])
	}

	n.Evaluation = val
	return val
}

func (e *EvalEngine) minmaxSerial(n *Node, depth int) float64 {
	if depth == 0 {
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	}

	if n.Children == nil {
		n.Children = n.GetChildNodes()
	}

	val := math.Inf(-1)
	var comp CompFunc = math.Max

	if n.Position.SideToMove == board.BlackToMove {
		val = math.Inf(1)
		comp = math.Min
	}

	for i := 0; i < len(n.Children); i++ {
		val = comp(val, e.minmax(n.Children[i], depth-1))
	}
	n.Evaluation = val
	return val
}
