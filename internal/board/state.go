package board

import (
	"math"
)

func (b *Board) updateCastlingRights(from Coord) {
	if b.CastlingRights == 0 {
		return
	}
	square := CoordToAlg(from)

	switch square {
	case "a1":
		b.CastlingRights = b.CastlingRights &^ WOOO
	case "h1":
		b.CastlingRights = b.CastlingRights &^ WOO
	case "a8":
		b.CastlingRights = b.CastlingRights &^ BOOO
	case "h8":
		b.CastlingRights = b.CastlingRights &^ BOO
	case "e1":
		b.CastlingRights = b.CastlingRights &^ WOOO
		b.CastlingRights = b.CastlingRights &^ WOO
	case "e8":
		b.CastlingRights = b.CastlingRights &^ BOOO
		b.CastlingRights = b.CastlingRights &^ BOO
	}
}

func (b *Board) updateEnPassantTarget(from, to Coord) {
	piece, _ := GetPiece(b, to)

	isPawnMove := piece == P || piece == p
	isDoubleMove := int(math.Abs(float64(from.Rank-to.Rank))) == 2

	if isDoubleMove && isPawnMove {
		b.EnPassantTarget = CoordToAlg(Coord{from.File, (from.Rank + to.Rank) / 2})
	} else {
		b.EnPassantTarget = "-"
	}
}

func (b *Board) updateSideToMove() {
	if b.SideToMove == WhiteToMove {
		b.SideToMove = BlackToMove
	} else {
		b.SideToMove = WhiteToMove
		b.FullMoveCounter++
	}
}
