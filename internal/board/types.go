package board

type Board struct {
	Coords          [8][8]uint8
	SideToMove      byte
	CastlingRights  CastlingRights
	EnPassantTarget string
	HalfMoveCounter uint8
	FullMoveCounter uint8
	IsEvaluated     bool
	CachedEval      float32
	TrackMoves      bool
	Moves           []string
}

type Coord struct {
	File int
	Rank int
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
