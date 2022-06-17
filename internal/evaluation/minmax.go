package eval

func (e *EvalEngine) minmaxSerial(n *Node, depth int, isWhite bool) float32 {
	if depth == 0 {
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	}

	if n.Children == nil {
		n.Children = n.GetChildNodes()
	}

	if len(n.Children) == 0 {
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	}

	var value float32
	if isWhite {
		value = negInf
		for i := 0; i < len(n.Children); i++ {
			value = Max32(value, e.minmaxSerial(n.Children[i], depth-1, false))
		}
		n.Evaluation = value
		return value

	} else {
		value = posInf
		for i := 0; i < len(n.Children); i++ {
			value = Min32(value, e.minmaxSerial(n.Children[i], depth-1, true))
		}
		n.Evaluation = value
		return value
	}

}
