package board

import (
	"os"
)

var CastlingMoves = [4]string{"e1g1", "e1c1", "e8g8", "e8c8"}
var Files = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}

func (b *Board) Init() {
	fen := os.Getenv("STARTING_FEN")
	if fen == "" {
		fen = startingFEN
	}
	b.ImportFEN(fen)
}

func (b *Board) MoveLongAlg(longalg string) {
	from, to := longAlgToCoords(longalg)
	if b.VerifyMove(longalg) {
		if b.trackMoves {
			b.TrackMove(longalg)
		}
		switch {
		case b.isCastling(longalg):
			b.castle(longalg)
		case b.isEnPassant(longalg):
			b.coords[to.File][to.Rank] = b.coords[from.File][from.Rank]
			b.coords[from.File][from.Rank] = empty
			b.coords[to.File][from.Rank] = empty
		default:
			b.coords[to.File][to.Rank] = b.coords[from.File][from.Rank]
			b.coords[from.File][from.Rank] = empty
		}

		b.updateEnPassantTarget(from, to)
		b.updateCastlingRights(from)
		b.autoPromotePawn(to)
		b.updateSideToMove()
	}
}

func (b *Board) castle(move string) {
	switch move {
	case "e1g1":
		b.coords[4][0] = empty
		b.coords[6][0] = K
		b.coords[7][0] = empty
		b.coords[5][0] = R
	case "e1c1":
		b.coords[4][0] = empty
		b.coords[2][0] = K
		b.coords[0][0] = empty
		b.coords[3][0] = R
	case "e8g8":
		b.coords[4][7] = empty
		b.coords[6][7] = k
		b.coords[7][7] = empty
		b.coords[5][7] = r
	case "e8c8":
		b.coords[4][7] = empty
		b.coords[2][7] = k
		b.coords[7][7] = empty
		b.coords[3][7] = r
	}
}

func (b *Board) autoPromotePawn(to Coord) {
	piece, _ := getPiece(*b, to)
	if piece == P && to.Rank == 7 {
		b.coords[to.File][to.Rank] = Q
	}

	if piece == p && to.Rank == 0 {
		b.coords[to.File][to.Rank] = q
	}
}

func (b Board) AccessCoord(c Coord) int {
	return b.coords[c.File][c.Rank]
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
	b.trackMoves = trackmoves
}

func (b *Board) TrackMove(move string) {
	//b.moves = append(b.moves, b.MoveToPretty(move))
	b.moves = append(b.moves, move)
}

func (b Board) GetMoveList() []string {
	return b.moves
}
