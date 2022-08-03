package eval

import "github.com/likeawizard/chess-go/internal/board"

type ttType uint8

const (
	TT_UPPER ttType = 1 << iota
	TT_LOWER
	TT_EXACT
)

type ttEntry struct {
	ttType ttType
	eval   int
	depth  int
	move   board.Move
}
