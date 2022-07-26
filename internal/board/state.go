package board

import (
	"math"
)

func (b *Board) updateCastlingRights(from, to Coord) {
	if b.CastlingRights == 0 {
		return
	}
	squareFrom := CoordToAlg(from)
	squareTo := CoordToAlg(to)

	switch {
	case b.CastlingRights&(WOOO|WOO) != 0 && squareFrom == "e1":
		b.ZobristCastlingRights(WOOO)
		b.ZobristCastlingRights(WOO)
		b.CastlingRights = b.CastlingRights &^ WOOO
		b.CastlingRights = b.CastlingRights &^ WOO

	case b.CastlingRights&(BOOO|BOO) != 0 && squareFrom == "e8":
		b.ZobristCastlingRights(BOOO)
		b.ZobristCastlingRights(BOO)
		b.CastlingRights = b.CastlingRights &^ BOOO
		b.CastlingRights = b.CastlingRights &^ BOO

	case b.CastlingRights&WOOO != 0 && (squareFrom == "a1" || squareTo == "a1"):
		b.ZobristCastlingRights(WOOO)
		b.CastlingRights = b.CastlingRights &^ WOOO

	case b.CastlingRights&WOO != 0 && (squareFrom == "h1" || squareTo == "h1"):
		b.ZobristCastlingRights(WOO)
		b.CastlingRights = b.CastlingRights &^ WOO

	case b.CastlingRights&BOOO != 0 && (squareFrom == "a8" || squareTo == "a8"):
		b.ZobristCastlingRights(BOOO)
		b.CastlingRights = b.CastlingRights &^ BOOO

	case b.CastlingRights&BOO != 0 && (squareFrom == "h8" || squareTo == "h8"):
		b.ZobristCastlingRights(BOO)
		b.CastlingRights = b.CastlingRights &^ BOO
	}
}

func (b *Board) updateEnPassantTarget(from, to Coord) {
	piece, _ := GetPiece(b, to)

	isPawnMove := piece == P || piece == p
	isDoubleMove := int(math.Abs(float64(from.Rank-to.Rank))) == 2

	if b.EnPassantTarget != "-" {
		b.ZobristEnPassant(b.EnPassantTarget)
	}

	if isDoubleMove && isPawnMove {
		b.EnPassantTarget = CoordToAlg(Coord{from.File, (from.Rank + to.Rank) / 2})
		b.ZobristEnPassant(b.EnPassantTarget)
	} else {
		b.EnPassantTarget = "-"
	}
}

func (b *Board) updateSideToMove() {
	b.ZobristSideToMove()
	if b.SideToMove == WhiteToMove {
		b.SideToMove = BlackToMove
	} else {
		b.SideToMove = WhiteToMove
		b.FullMoveCounter++
	}
}
