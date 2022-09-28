package eval

import "github.com/likeawizard/chess-go/internal/board"

// Piece-square-tables used for positional evaluation of piece placement
// The tables are asymmetrical and are viewed from white's perspective.

func init() {
	invert := func(sq int) int {
		return (7-sq/8)*8 + sq%8
	}
	for piece := board.PAWNS; piece <= board.KINGS; piece++ {
		for sq := 0; sq < 64; sq++ {
			PST[board.BLACK][piece][sq] = PST[board.WHITE][piece][invert(sq)]
		}
	}
}

var PST = [2][6][64]int{{pawnPST, bishopPST, kingPST, rookPST, queenPST, kingPST}}

var pawnPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 10, 10, 10, 10, 10, 10,
	7, 7, 7, 8, 8, 7, 7, 7,
	5, 5, 5, 10, 10, 5, 5, 5,
	4, 5, 4, 15, 15, 10, 5, 4,
	3, 4, 0, 3, 3, -5, 4, 3,
	2, 5, 5, 0, 0, 5, 5, 2,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var bishopPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 4, 0, 0, 1, 0, 0,
	0, 5, 0, 5, 1, 0, 3, 0,
	0, 0, 2, 1, 5, 3, 0, 0,
	0, 0, 1, 0, 0, 4, 0, 0,
	0, 10, 0, 0, 0, 0, 10, 0,
	0, 0, -2, 0, 0, -2, 0, 0,
}

var knightPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 5, 5, 5, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 2, 3, 5, 5, 3, 2, 0,
	0, 0, 0, 1, 1, 0, 0, 0,
	0, -10, 0, 0, 0, 0, -10, 0,
}

var rookPST [64]int = [64]int{
	0, 0, 2, 3, 3, 0, 0, 0,
	5, 5, 5, 5, 5, 5, 5, 5,
	0, 0, 2, 3, 3, 2, 0, 0,
	0, 0, 2, 3, 3, 2, 0, 0,
	0, 0, 2, 3, 3, 2, 0, 0,
	0, 0, 2, 3, 3, 2, 0, 0,
	0, 0, 2, 3, 3, 2, 0, 0,
	0, 0, 2, 3, 3, 2, 0, 0,
}

var queenPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 2, 0, 0, 0, 0, 0, 0,
	0, 0, 2, 0, 0, 0, 0, 0,
	0, 0, 0, 3, 0, 0, 0, 0,
}

var kingPST [64]int = [64]int{
	-5, -8, -10, -10, -10, -10, -8, -5,
	-5, -8, -10, -10, -10, -10, -8, -5,
	-5, -8, -10, -10, -10, -10, -8, -5,
	-5, -8, -10, -10, -10, -10, -8, -5,
	-5, -8, -10, -10, -10, -10, -8, -5,
	-3, -3, -4, -5, -5, -4, -4, -3,
	-1, -1, -2, -2, -2, -2, -1, -1,
	0, 0, 10, 0, 0, 0, 10, 0,
}
