package eval

import "github.com/likeawizard/chess-go/internal/board"

// Piece-square-tables used for positional evaluation of piece placement
// The tables are asymmetrical and are viewed from white's perspective.

func init() {
	invert := func(sq int) int {
		return (7-sq/8)*8 + sq%8
	}
	for stage := 0; stage < 2; stage++ {
		for piece := board.PAWNS; piece <= board.KINGS; piece++ {
			for sq := 0; sq < 64; sq++ {
				PST[stage][board.BLACK][piece][sq] = PST[stage][board.WHITE][piece][invert(sq)]
			}
		}
	}
}

var PST = [2][2][6][64]int{
	{{pawnPST, bishopPST, knightPST, rookPST, queenPST, kingPST}},
	{{pawnEGPST, bishopEGPST, knightEGPST, rookEGPST, queenEGPST, kingEGPST}},
}

var pawnPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 10, 10, 10, 10, 10, 10,
	7, 7, 7, 8, 8, 7, 7, 7,
	5, 5, 15, 20, 20, 5, 5, 5,
	4, 5, 20, 35, 35, 7, 5, 4,
	3, 4, 0, 10, 10, -5, 4, 3,
	2, 5, 5, -15, -15, 5, 5, 2,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var pawnEGPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	30, 30, 30, 30, 30, 30, 30, 30,
	20, 14, 14, 8, 8, 14, 14, 20,
	16, 12, 10, 6, 6, 10, 12, 16,
	12, 10, 8, 4, 4, 8, 10, 12,
	8, 8, 4, 4, 4, 4, 8, 8,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var bishopPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 8, 0, 0, 8, 0, 0,
	0, 10, 0, 10, 10, 0, 6, 0,
	0, 0, 2, 10, 10, 6, 0, 0,
	0, 0, 8, 0, 0, 8, 0, 0,
	0, 20, 0, 0, 0, 0, 10, 0,
	0, 0, -10, 0, 0, -10, 0, 0,
}

var bishopEGPST [64]int = [64]int{
	10, 10, 0, 0, 0, 0, 10, 10,
	10, 20, 20, 10, 10, 20, 20, 10,
	0, 20, 40, 30, 30, 40, 20, 0,
	0, 10, 30, 60, 60, 30, 10, 0,
	0, 10, 30, 60, 60, 30, 10, 0,
	0, 20, 40, 30, 30, 40, 20, 0,
	10, 20, 20, 10, 10, 20, 20, 10,
	10, 10, 0, 0, 0, 0, 10, 10,
}

var knightPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 10, 10, 10, 10, 0, 0,
	0, 0, 10, 10, 10, 10, 0, 0,
	0, 0, 10, 10, 10, 10, 0, 0,
	0, 4, 6, 10, 10, 6, 4, 0,
	0, 0, 0, 3, 3, 0, 0, 0,
	0, -10, 0, 0, 0, 0, -10, 0,
}

var knightEGPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 3, 6, 6, 6, 6, 3, 0,
	0, 2, 12, 14, 14, 12, 6, 0,
	0, 2, 14, 35, 35, 14, 6, 0,
	0, 2, 14, 35, 35, 14, 6, 0,
	0, 2, 12, 14, 14, 12, 6, 0,
	0, 3, 6, 6, 6, 6, 3, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var rookPST [64]int = [64]int{
	0, 0, 8, 9, 9, 0, 0, 0,
	15, 15, 15, 15, 15, 15, 15, 15,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
}

var rookEGPST [64]int = [64]int{
	10, 10, 10, 15, 15, 10, 10, 10,
	15, 15, 15, 15, 15, 15, 15, 15,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
	0, 0, 8, 9, 9, 8, 0, 0,
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

var queenEGPST [64]int = [64]int{
	1, 1, 0, 0, 0, 0, 1, 1,
	1, 2, 2, 1, 1, 2, 2, 1,
	0, 2, 6, 5, 5, 6, 2, 0,
	0, 1, 5, 10, 10, 5, 1, 0,
	0, 1, 5, 10, 10, 5, 1, 0,
	0, 2, 6, 5, 5, 6, 2, 0,
	1, 2, 2, 1, 1, 2, 2, 1,
	1, 1, 0, 0, 0, 0, 1, 1,
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

var kingEGPST [64]int = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 3, 2, 2, 2, 2, 3, 0,
	0, 4, 6, 7, 7, 6, 4, 0,
	0, 2, 7, 10, 10, 7, 2, 0,
	0, 2, 7, 10, 10, 7, 2, 0,
	0, 2, 6, 7, 7, 6, 2, 0,
	0, 1, 2, 2, 2, 2, 1, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}
