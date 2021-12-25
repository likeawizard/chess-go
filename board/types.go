package board

/*
0 empty
1 white pawn
2 wb
3 wk
4 wr
5 wq
6 wk
7 bp
8 bb
9 bk
10 br
11 bq
12 bk
*/
type Board struct {
	coords          [8][8]int
	sideToMove      string
	castlingRights  string
	enPassantTarget string
	halfMoveCounter int
	fullMoveCounter int
	isEvaluated     bool
	cachedEval      float64
	enPassantMoves  []string
	trackMoves      bool
	moves           []string
}

type Coord struct {
	File int
	Rank int
}

const (
	empty = 0
	P     = 1
	B     = 2
	N     = 3
	R     = 4
	Q     = 5
	K     = 6
	p     = 7
	b     = 8
	n     = 9
	r     = 10
	q     = 11
	k     = 12
)

const (
	startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	whiteToMove = "w"
	blackToMove = "b"
	wOO         = "K"
	wOOO        = "Q"
	bOO         = "k"
	bOOO        = "q"
)
