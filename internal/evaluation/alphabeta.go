package eval

import (
	"math"
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

func (e *EvalEngine) alphabeta(n *Node, depth int, alpha, beta float64) float64 {
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

		if e.SearchDepth-depth < 1 {
			wg.Add(1)
			e.MaxGoroutines <- struct{}{}
			go func(i int) {
				defer wg.Done()

				temp[i] = e.alphabeta(n.Children[i], depth-1, alpha, beta)
				<-e.MaxGoroutines
			}(i)
		} else {
			temp[i] = e.alphabetaSerial(n.Children[i], depth-1, alpha, beta)
		}

	}

	wg.Wait()

	for i := 0; i < len(temp); i++ {
		val = comp(val, temp[i])
	}

	n.Evaluation = val
	return val

}

func (e *EvalEngine) alphabetaSerial(n *Node, depth int, alpha, beta float64) float64 {
	if depth == 0 {
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	}

	if n.Children == nil {
		n.Children = n.GetChildNodes()
	}

	val := math.Inf(-1)
	var comp CompFunc = math.Max
	var compB CompFuncBool = gte
	var selectivecomp SelectiveCompFunc = maxA

	if n.Position.SideToMove == board.BlackToMove {
		val = math.Inf(1)
		comp = math.Min
		compB = lte
		selectivecomp = minB
	}

	for i := 0; i < len(n.Children); i++ {
		val = comp(val, e.alphabetaSerial(n.Children[i], depth-1, alpha, beta))
		alpha, beta = selectivecomp(val, alpha, beta)
		if compB(val, alpha, beta) {
			break
		}
	}
	n.Evaluation = val
	return val
}
