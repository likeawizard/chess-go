package board

import (
	"math"
	"strings"
)

func (b *Board) removeCastlingRight(castling string) {
	b.castlingRights = strings.ReplaceAll(b.castlingRights, castling, "")
	if len(b.castlingRights) == 0 {
		b.castlingRights = "-"
	}
}

func (b *Board) updateCastlingRights(from Coord) {
	if b.castlingRights == "-" {
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
	piece, _ := getPiece(*b, to)

	isPawnMove := piece == P || piece == p
	isDoubleMove := int(math.Abs(float64(from.Rank-to.Rank))) == 2

	if isDoubleMove && isPawnMove {
		b.enPassantTarget = CoordToAlg(Coord{from.File, (from.Rank + to.Rank) / 2})
	} else {
		b.enPassantTarget = "-"
	}
}

func (b *Board) updateSideToMove() {
	if b.sideToMove == whiteToMove {
		b.sideToMove = blackToMove
	} else {
		b.sideToMove = whiteToMove
		b.fullMoveCounter++
	}
}
