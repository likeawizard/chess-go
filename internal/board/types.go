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
	Coords          [8][8]uint8
	SideToMove      byte
	CastlingRights  string
	EnPassantTarget string
	HalfMoveCounter uint8
	FullMoveCounter uint8
	IsEvaluated     bool
	CachedEval      float32
	EnPassantMoves  []string
	TrackMoves      bool
	Moves           []string
}

type Coord struct {
	File int
	Rank int
}

const (
	empty uint8 = 0
	P     uint8 = 1
	B     uint8 = 2
	N     uint8 = 3
	R     uint8 = 4
	Q     uint8 = 5
	K     uint8 = 6
	p     uint8 = 7
	b     uint8 = 8
	n     uint8 = 9
	r     uint8 = 10
	q     uint8 = 11
	k     uint8 = 12
)

const (
	WhiteToMove byte = 'w'
	BlackToMove byte = 'b'
)

const (
	startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	wOO         = "K"
	wOOO        = "Q"
	bOO         = "k"
	bOOO        = "q"
)
