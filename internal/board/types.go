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
	Hash            uint64
	Coords          [64]uint8
	SideToMove      byte
	CastlingRights  CastlingRights
	EnPassantTarget Square
	HalfMoveCounter uint8
	FullMoveCounter uint8
	IsEvaluated     bool
	CachedEval      float32
	TrackMoves      bool
	Moves           []Move
}

const (
	empty uint8 = iota
	P
	B
	N
	R
	Q
	K
	p
	b
	n
	r
	q
	k
)

const (
	WhiteToMove byte = 'w'
	BlackToMove byte = 'b'
)

type CastlingRights byte

const (
	WOO CastlingRights = 1 << iota
	WOOO
	BOO
	BOOO
	CASTLING_ALL = WOO | WOOO | BOO | BOOO
)

const (
	startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)
