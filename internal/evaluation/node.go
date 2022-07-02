package eval

import (
	"context"
	"sort"

	"github.com/likeawizard/chess-go/internal/board"
)

type Node struct {
	Position   *board.Board
	MoveToPlay string
	Evaluation float32
	Parent     *Node
	Children   []*Node
}

type SearchFunction func(n *Node, depth ...int) float32

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

func (n *Node) BuildGameTree(depth int) {
	if depth == 0 {
		return
	}
	n.Children = n.GetChildNodes()
	for _, child := range n.Children {
		child.BuildGameTree(depth - 1)
	}
}

func (n *Node) EvaluateLeafNodes(ctx context.Context, e *EvalEngine) {
	select {
	case <-ctx.Done():
		return
	default:
		if n.Children == nil || len(n.Children) == 0 {
			GetEvaluation(e, n.Position)
		} else {
			for _, child := range n.Children {
				child.EvaluateLeafNodes(ctx, e)
			}
		}
	}
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
			if c.Evaluation >= bestScore {
				bestScore, bestMove = c.Evaluation, c
			}
		}
	case board.BlackToMove:
		bestScore = posInf
		for _, c := range n.Children {
			if c.Evaluation <= bestScore {
				bestScore, bestMove = c.Evaluation, c
			}
		}
	}

	return bestMove
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
	line = append(line, n.MoveToPlay)
	side := n.Position.SideToMove
	current := n
	for current.Children != nil {
		best := current.PickBestMove(side)
		if best == nil {
			break
		}
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

// traverse(Node node) {
//     if (node==NULL)
//         return;

//     stack<Node> stk;
//     stk.push(node);

//     while (!stk.empty()) {
//         Node top = stk.pop();
//         for (Node child in top.getChildren()) {
//             stk.push(child);
//         }
//         process(top);
//     }
// }
func (n *Node) EvaluateLeafNodesNR(e *EvalEngine) {
	stack := NewStack()

	stack.Push(n)
	for !stack.Empty() {
		top := stack.Pop()
		if top.Children != nil {
			stack.PushSlice(top.Children)
		} else {
			GetEvaluation(e, top.Position)
		}
	}
}

type NodeStack []*Node

func NewStack() NodeStack {
	return make([]*Node, 0)
}

func (ns *NodeStack) Push(n *Node) {
	*ns = append(*ns, n)
}

func (ns *NodeStack) PushSlice(n []*Node) {
	*ns = append(*ns, n...)
}

func (ns *NodeStack) Pop() *Node {
	if ns.Empty() {
		return nil
	} else {
		last := len(*ns) - 1
		top := (*ns)[last]
		*ns = (*ns)[:last]
		return top
	}
}

func (ns *NodeStack) Empty() bool {
	return ns == nil || len(*ns) == 0
}
