package board

import (
	"math/rand"
)

func init() {
	seed = rand.Uint64()
	castlingKeys = make(map[CastlingRights]uint64)
	for sq := 0; sq < 64; sq++ {
		for i := 0; i <= 12; i++ {
			pieceKeys[i][sq] = rand.Uint64()
		}
		enPassantKeys[sq] = rand.Uint64()
	}

	castlingKeys[WOO] = rand.Uint64()
	castlingKeys[WOOO] = rand.Uint64()
	castlingKeys[BOO] = rand.Uint64()
	castlingKeys[BOOO] = rand.Uint64()

	swapSide = rand.Uint64()
}

var seed uint64
var pieceKeys [14][64]uint64
var castlingKeys map[CastlingRights]uint64
var swapSide uint64
var enPassantKeys [64]uint64

func (b *Board) SeedHash() uint64 {
	hash := seed

	var p int
	for sq := 0; sq < 64; sq++ {
		p = int(b.Coords[sq])
		hash ^= pieceKeys[p][sq]

	}

	for _, cr := range castlingKeys {
		hash ^= cr
	}

	if b.EnPassantTarget != -1 {
		hash ^= enPassantKeys[b.EnPassantTarget]
	}

	return hash
}

func (b *Board) ZobristSimpleMove(move Move) {
	from, to := move.From(), move.To()
	start := b.Coords[from]
	finish := b.Coords[to]

	// unset target piece at destination and set new
	b.Hash ^= pieceKeys[finish][to]
	b.Hash ^= pieceKeys[start][to]
	// unset moved piece and replace with empty
	b.Hash ^= pieceKeys[start][from]
	b.Hash ^= pieceKeys[0][from]
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
		b.ZobristSimpleMove(WCastleKing)
		b.ZobristSimpleMove(WCastleKingRook)
	case WOOO:
		b.ZobristSimpleMove(WCastleQueen)
		b.ZobristSimpleMove(WCastleQueenRook)
	case BOO:
		b.ZobristSimpleMove(BCastleKing)
		b.ZobristSimpleMove(BCastleKingRook)
	case BOOO:
		b.ZobristSimpleMove(BCastleQueen)
		b.ZobristSimpleMove(BCastleQueenRook)
	}
}

func (b *Board) ZobristPromotion(move Move) {
	offset := uint8(6)
	if b.IsWhite {
		offset = 0
	}
	promotion := move.Promotion()
	switch promotion {
	case 'q':
		promotion = Q + offset
	case 'n':
		promotion = N + offset
	case 'r':
		promotion = R + offset
	case 'b':
		promotion = B + offset
	}
	from, to := move.From(), move.To()
	start := b.Coords[from]
	finish := b.Coords[to]

	b.Hash ^= pieceKeys[finish][to]
	// set destination with newly promoted piece
	b.Hash ^= pieceKeys[promotion][to]
	b.Hash ^= pieceKeys[finish][from]
	b.Hash ^= pieceKeys[start][from]

}

func (b *Board) ZobristEnPassant(square Square) {
	if b.EnPassantTarget != -1 {
		b.Hash ^= enPassantKeys[square]
	}
}
