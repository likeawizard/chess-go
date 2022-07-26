package board

import (
	"fmt"
	"time"
)

func Perft(fen string, depth int) (int, time.Duration) {
	b := &Board{}
	b.ImportFEN(fen)
	root := Node{Position: b}
	start := time.Now()
	root.BuildGameTree(depth)
	leafs := root.CountLeafs()
	return leafs, time.Since(start)
}

func PerftDebug(fen string, depth int) {
	b := &Board{}
	b.ImportFEN(fen)
	root := Node{Position: b}
	root.BuildGameTree(depth)
	for _, child := range root.Children {
		fmt.Printf("%s: %d\n", child.MoveToPlay, child.CountLeafs())
	}
}

type Node struct {
	Position   *Board
	MoveToPlay string
	Parent     *Node
	Children   []*Node
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

func (n *Node) GetChildNodes() []*Node {
	moves, captures := n.Position.GetLegalMoves(n.Position.SideToMove)
	all := append(captures, moves...)
	childNodes := make([]*Node, len(all))

	for i := 0; i < len(all); i++ {
		childNodes[i] = &Node{
			Parent:     n,
			Children:   nil,
			Position:   &Board{},
			MoveToPlay: all[i],
		}
		childNodes[i].Position = n.Position.SimpleCopy()
		childNodes[i].Position.MoveLongAlg(all[i])
	}

	return childNodes
}

func (n *Node) CountLeafs() int {
	num := 0
	if n.Children == nil {
		return 1
	} else {
		for _, c := range n.Children {
			num += c.CountLeafs()
		}
	}

	return num
}
