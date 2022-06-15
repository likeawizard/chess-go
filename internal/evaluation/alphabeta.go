package eval

import (
	"sync"

	"github.com/likeawizard/chess-go/internal/board"
)

func (e *EvalEngine) alphabeta(n *Node, depth int, alpha, beta float32) float32 {
	if depth == 0 {
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	}

	if n.Children == nil {
		n.Children = n.GetChildNodes()
	}

	val := negInf
	var comp CompFunc = Max32

	if n.Position.SideToMove == board.BlackToMove {
		val = posInf
		comp = Min32
	}

	var wg sync.WaitGroup
	temp := make([]float32, len(n.Children))

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

func (e *EvalEngine) alphabetaSerial(n *Node, depth int, alpha, beta float32) float32 {
	if depth == 0 {
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	}

	if n.Children == nil {
		n.Children = n.GetChildNodes()
	}

	val := negInf
	var comp CompFunc = Max32
	var compB CompFuncBool = gte
	var selectivecomp SelectiveCompFunc = maxA

	if n.Position.SideToMove == board.BlackToMove {
		val = posInf
		comp = Min32
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
