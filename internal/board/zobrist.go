package board

import (
	"math/rand"
)

func init() {
	seed = rand.Uint64()
	castlingKeys = make(map[CastlingRights]uint64)
	for sq := 0; sq < 64; sq++ {
		for color := WHITE; color <= BLACK; color++ {
			for pieceType := PAWNS; pieceType <= KINGS; pieceType++ {
				pieceKeys[color][pieceType][sq] = rand.Uint64()
			}
		}
		for i := 1; i <= 12; i++ {

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
var pieceKeys [2][6][64]uint64
var castlingKeys map[CastlingRights]uint64
var swapSide uint64
var enPassantKeys [64]uint64

// Calculate Zborist hash of the position
func (b *Board) SeedHash() uint64 {
	hash := seed

	for color := WHITE; color <= BLACK; color++ {
		for pieceType := PAWNS; pieceType <= KINGS; pieceType++ {
			pieces := b.Pieces[color][pieceType]
			for pieces > 0 {
				sq := pieces.PopLS1B()
				hash ^= pieceKeys[color][pieceType][sq]
			}
		}
	}

	for _, cr := range castlingKeys {
		hash ^= cr
	}

	if b.EnPassantTarget != -1 {
		hash ^= enPassantKeys[b.EnPassantTarget]
	}

	return hash
}

// Incrementally update Zborist hash after a move
// TODO: optimize - remove use of expensive PieceAtSquare function
func (b *Board) ZobristSimpleMove(move Move) {
	from, to := move.From(), move.To()
	_, color, piece := b.PieceAtSquare(from)

	// unset target piece at destination and set new
	if move.IsCapture() {
		_, _, capturedPiece := b.PieceAtSquare(to)
		b.Hash ^= pieceKeys[color^1][capturedPiece][to]
	}

	b.Hash ^= pieceKeys[color][piece][to]
	b.Hash ^= pieceKeys[color][piece][from]
}

// Update Zobirst hash with flipping side to move
func (b *Board) ZobristSideToMove() {
	b.Hash ^= swapSide
}

// Update Zobrist hash with castling rights
func (b *Board) ZobristCastlingRights(right CastlingRights) {
	b.Hash ^= castlingKeys[right]
}

// Update Zobrist hash with castling move
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

// Update Zobrist hash when promoting a piece
func (b *Board) ZobristPromotion(move Move) {
	var promotion int
	switch move.Promotion() {
	case 'q':
		promotion = QUEENS
	case 'n':
		promotion = KNIGHTS
	case 'r':
		promotion = ROOKS
	case 'b':
		promotion = BISHOPS
	}
	from, to := move.From(), move.To()
	_, color, _ := b.PieceAtSquare(from)

	if move.IsCapture() {
		_, _, capturedPiece := b.PieceAtSquare(to)
		b.Hash ^= pieceKeys[color^1][capturedPiece][to]
	}

	// set destination with newly promoted piece
	b.Hash ^= pieceKeys[color][promotion][to]
	b.Hash ^= pieceKeys[color][PAWNS][from]
}

// Update Zobrist hash with En Passant square
func (b *Board) ZobristEnPassant(square Square) {
	if b.EnPassantTarget != -1 {
		b.Hash ^= enPassantKeys[square]
	}
}
