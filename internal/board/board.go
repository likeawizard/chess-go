package board

import (
	"strings"

	"github.com/likeawizard/chess-go/internal/config"
)

var CastlingMoves = [4]string{"e1g1", "e1c1", "e8g8", "e8c8"}
var CastlingRights = [4]string{wOO, wOOO, bOO, bOOO}
var Files = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}

func (b *Board) Init(c *config.Config) {
	fen := c.Init.StartingFen
	if fen == "" {
		fen = startingFEN
	}
	b.ImportFEN(fen)
}

func (b *Board) InitDefault() {
	b.ImportFEN(startingFEN)
}

func (b *Board) Copy() *Board {
	return &Board{
		Coords:          b.Coords,
		SideToMove:      b.SideToMove,
		CastlingRights:  b.CastlingRights,
		EnPassantTarget: b.EnPassantTarget,
		HalfMoveCounter: b.HalfMoveCounter,
		FullMoveCounter: b.HalfMoveCounter,
		IsEvaluated:     b.IsEvaluated,
		CachedEval:      b.CachedEval,
		EnPassantMoves:  b.EnPassantMoves,
		TrackMoves:      b.TrackMoves,
		Moves:           b.Moves,
	}
}

// Only copy fields necessary for gametree construction
func (b *Board) SimpleCopy() *Board {
	return &Board{
		Coords:          b.Coords,
		SideToMove:      b.SideToMove,
		CastlingRights:  b.CastlingRights,
		EnPassantTarget: b.EnPassantTarget,
		HalfMoveCounter: b.HalfMoveCounter,
		FullMoveCounter: b.HalfMoveCounter,
	}
}

func (b *Board) MoveLongAlg(longalg string) {
	from, to := longAlgToCoords(longalg)
	if b.VerifyMove(longalg) {
		if b.TrackMoves {
			b.TrackMove(longalg)
		}
		switch {
		case b.isCastling(longalg):
			b.castle(longalg)
		case b.isEnPassant(longalg):
			b.Coords[to.File][to.Rank] = b.Coords[from.File][from.Rank]
			b.Coords[from.File][from.Rank] = empty
			b.Coords[to.File][from.Rank] = empty
		case len(longalg) == 5:
			b.Coords[to.File][to.Rank] = b.promote(longalg[4:])
			b.Coords[from.File][from.Rank] = empty
		default:
			b.Coords[to.File][to.Rank] = b.Coords[from.File][from.Rank]
			b.Coords[from.File][from.Rank] = empty
		}

		b.updateEnPassantTarget(from, to)
		b.updateCastlingRights(from)
		b.updateSideToMove()
	}
}

func (b *Board) promote(piece string) uint8 {
	off := uint8(0)
	if b.SideToMove == BlackToMove {
		off = 6
	}

	switch piece {
	case "b":
		return off + B
	case "r":
		return off + R
	case "n":
		return off + N
	default:
		return off + Q
	}
}

func (b *Board) castle(move string) {
	switch move {
	case "e1g1":
		b.Coords[4][0] = empty
		b.Coords[6][0] = K
		b.Coords[7][0] = empty
		b.Coords[5][0] = R
	case "e1c1":
		b.Coords[4][0] = empty
		b.Coords[2][0] = K
		b.Coords[0][0] = empty
		b.Coords[3][0] = R
	case "e8g8":
		b.Coords[4][7] = empty
		b.Coords[6][7] = k
		b.Coords[7][7] = empty
		b.Coords[5][7] = r
	case "e8c8":
		b.Coords[4][7] = empty
		b.Coords[2][7] = k
		b.Coords[0][7] = empty
		b.Coords[3][7] = r
	}
}

func (b *Board) AccessCoord(c Coord) uint8 {
	return b.Coords[c.File][c.Rank]
}

func CoordInBounds(c Coord) bool {
	return c.Rank <= 7 && c.Rank >= 0 && c.File <= 7 && c.File >= 0
}

func longAlgToCoords(longalg string) (from, to Coord) {
	from = AlgToCoord(longalg[:2])
	to = AlgToCoord(longalg[2:])

	return
}

func fileToCoord(file rune) int {
	for i, f := range Files {
		if f == string(file) {
			return i
		}
	}
	return 0
}

func (c *Coord) Equal(a *Coord) bool {
	return c.File == a.File && c.Rank == a.Rank
}

func CoordToAlg(c Coord) string {
	return Files[c.File] + string(rune(c.Rank+1+'0'))
}

func AlgToCoord(alg string) (c Coord) {
	chars := []rune(alg)
	c = Coord{File: fileToCoord(chars[0]), Rank: int(chars[1]-'0') - 1}
	return c
}

func CoordsToMove(from, to Coord) string {
	return CoordToAlg(from) + CoordToAlg(to)
}

func (b *Board) SetTrackMoves(trackmoves bool) {
	b.TrackMoves = trackmoves
}

func (b *Board) TrackMove(move string) {
	b.Moves = append(b.Moves, move)
}

func (b *Board) GetMoveList() []string {
	return b.Moves
}

func (b *Board) GetLastMove() string {
	if len(b.Moves) == 0 {
		return ""
	}
	return b.Moves[len(b.Moves)-1]
}

func (b *Board) PlayMoves(moves string) {
	moveSlice := strings.Fields(moves)
	for _, move := range moveSlice {
		b.MoveLongAlg(move)
	}
}
