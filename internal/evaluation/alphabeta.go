package eval

func (e *EvalEngine) alphabetaSerial(n *Node, depth int, alpha, beta float32, isWhite bool) float32 {
	if depth == 0 {
		n.Evaluation = e.EvalFunction(e, n.Position)
		return n.Evaluation
	}

	if n.Children == nil {
		n.Children = n.GetChildNodes()
	}

	var value float32
	if isWhite {
		value = negInf
		for i := 0; i < len(n.Children); i++ {

			value = Max32(value, e.alphabetaSerial(n.Children[i], depth-1, alpha, beta, false))

			if value >= beta {
				break
			}
			alpha = Max32(alpha, value)
		}
		if n != nil {
			n.Evaluation = value
		}

		return value
	} else {
		value = posInf
		for i := 0; i < len(n.Children); i++ {
			value = Min32(value, e.alphabetaSerial(n.Children[i], depth-1, alpha, beta, true))

			if value <= alpha {
				break
			}
			beta = Min32(beta, value)
		}
		if n != nil {
			n.Evaluation = value
		}
		return value
	}
}
