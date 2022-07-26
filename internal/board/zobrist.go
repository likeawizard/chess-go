package board

import (
	"math/rand"
)

func init() {
	seed = rand.Uint64()
	castlingKeys = make(map[CastlingRights]uint64)
	enPassantKeys = make(map[string]uint64)
	for i := 0; i <= 12; i++ {
		for f := 0; f < 8; f++ {
			for r := 0; r < 8; r++ {
				pieceKeys[i][f][r] = rand.Uint64()
			}
		}
	}
	castlingKeys[WOO] = rand.Uint64()
	castlingKeys[WOOO] = rand.Uint64()
	castlingKeys[BOO] = rand.Uint64()
	castlingKeys[BOOO] = rand.Uint64()

	swapSide = rand.Uint64()
}

var seed uint64
var pieceKeys [14][8][8]uint64
var castlingKeys map[CastlingRights]uint64
var swapSide uint64
var enPassantKeys map[string]uint64

func (b *Board) SeedHash() uint64 {
	hash := seed

	var p int
	for f := 0; f < 8; f++ {
		for r := 0; r < 8; r++ {
			p = int(b.Coords[f][r])
			hash ^= pieceKeys[p][f][r]
		}
	}

	for _, cr := range castlingKeys {
		hash ^= cr
	}

	var c Coord
	var square string
	for f := 0; f < 8; f++ {
		for r := 0; r < 8; r++ {
			c.File, c.Rank = f, r
			square = CoordToAlg(c)
			enPassantKeys[square] = rand.Uint64()
		}
	}

	return hash
}

func (b *Board) ZobristSimpleMove(from, to Coord) {
	start := b.AccessCoord(from)
	finish := b.AccessCoord(to)

	// unset target piece at destination and set new
	b.Hash ^= pieceKeys[finish][to.File][to.Rank]
	b.Hash ^= pieceKeys[start][to.File][to.Rank]
	// unset moved piece and replace with empty
	b.Hash ^= pieceKeys[start][from.File][from.Rank]
	b.Hash ^= pieceKeys[0][from.File][from.Rank]
}

func (b *Board) ZobristSideToMove() {
	b.Hash ^= swapSide
}

func (b *Board) ZobristCastlingRights(right CastlingRights) {
	b.Hash ^= castlingKeys[right]
}

func (b *Board) ZobristCastling(right CastlingRights) {
	switch right {
	case WOO:
		kfrom := Coord{File: 4, Rank: 0}
		kto := Coord{File: 6, Rank: 0}
		rfrom := Coord{File: 7, Rank: 0}
		rto := Coord{File: 5, Rank: 0}
		b.ZobristSimpleMove(kfrom, kto)
		b.ZobristSimpleMove(rfrom, rto)
	case WOOO:
		kfrom := Coord{File: 4, Rank: 0}
		kto := Coord{File: 2, Rank: 0}
		rfrom := Coord{File: 0, Rank: 0}
		rto := Coord{File: 3, Rank: 0}
		b.ZobristSimpleMove(kfrom, kto)
		b.ZobristSimpleMove(rfrom, rto)
	case BOO:
		kfrom := Coord{File: 4, Rank: 7}
		kto := Coord{File: 6, Rank: 7}
		rfrom := Coord{File: 7, Rank: 7}
		rto := Coord{File: 5, Rank: 7}
		b.ZobristSimpleMove(kfrom, kto)
		b.ZobristSimpleMove(rfrom, rto)
	case BOOO:
		kfrom := Coord{File: 4, Rank: 7}
		kto := Coord{File: 2, Rank: 7}
		rfrom := Coord{File: 0, Rank: 7}
		rto := Coord{File: 3, Rank: 7}
		b.ZobristSimpleMove(kfrom, kto)
		b.ZobristSimpleMove(rfrom, rto)
	}
}

func (b *Board) ZobristPromotion(from, to Coord, promoteTo uint8) {
	promotion := int(promoteTo)
	start := b.AccessCoord(from)
	finish := b.AccessCoord(to)

	b.Hash ^= pieceKeys[finish][to.File][to.Rank]
	// set destination with newly promoted piece
	b.Hash ^= pieceKeys[promotion][to.File][to.Rank]
	b.Hash ^= pieceKeys[finish][from.File][from.Rank]
	b.Hash ^= pieceKeys[start][from.File][from.Rank]

}

func (b *Board) ZobristEnPassant(square string) {
	b.Hash ^= enPassantKeys[square]
}
