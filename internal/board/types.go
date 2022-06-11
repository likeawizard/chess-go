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
	Coords          [8][8]int
	SideToMove      string
	CastlingRights  string
	EnPassantTarget string
	HalfMoveCounter int
	FullMoveCounter int
	IsEvaluated     bool
	CachedEval      float64
	EnPassantMoves  []string
	TrackMoves      bool
	Moves           []string
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
	WhiteToMove = "w"
	BlackToMove = "b"
	wOO         = "K"
	wOOO        = "Q"
	bOO         = "k"
	bOOO        = "q"
)
