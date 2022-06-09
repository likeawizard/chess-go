package board

import (
	"math"
	"strings"
)

func (b *Board) removeCastlingRight(castling string) {
	b.CastlingRights = strings.ReplaceAll(b.CastlingRights, castling, "")
	if len(b.CastlingRights) == 0 {
		b.CastlingRights = "-"
	}
}

func (b *Board) updateCastlingRights(from Coord) {
	if b.CastlingRights == "-" {
		return
	}
	square := CoordToAlg(from)

	switch square {
	case "a1":
		b.removeCastlingRight(wOOO)
	case "h1":
		b.removeCastlingRight(wOO)
	case "a8":
		b.removeCastlingRight(bOOO)
	case "h8":
		b.removeCastlingRight(bOO)
	case "e1":
		b.removeCastlingRight(wOOO)
		b.removeCastlingRight(wOO)
	case "e8":
		b.removeCastlingRight(bOOO)
		b.removeCastlingRight(bOO)
	}
}

func (b *Board) updateEnPassantTarget(from, to Coord) {
	piece, _ := GetPiece(*b, to)

	isPawnMove := piece == P || piece == p
	isDoubleMove := int(math.Abs(float64(from.Rank-to.Rank))) == 2

	if isDoubleMove && isPawnMove {
		b.EnPassantTarget = CoordToAlg(Coord{from.File, (from.Rank + to.Rank) / 2})
	} else {
		b.EnPassantTarget = "-"
	}
}

func (b *Board) updateSideToMove() {
	if b.SideToMove == whiteToMove {
		b.SideToMove = blackToMove
	} else {
		b.SideToMove = whiteToMove
		b.FullMoveCounter++
	}
}
