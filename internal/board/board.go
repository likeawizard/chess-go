package board

import (
	"strings"

	"github.com/likeawizard/chess-go/internal/config"
)

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
		Hash:            b.Hash,
		Coords:          b.Coords,
		SideToMove:      b.SideToMove,
		CastlingRights:  b.CastlingRights,
		EnPassantTarget: b.EnPassantTarget,
		HalfMoveCounter: b.HalfMoveCounter,
		FullMoveCounter: b.HalfMoveCounter,
	}
}

func (b *Board) MoveLongAlg(move Move) {
	from, to := move.FromTo()
	switch {
	case b.IsCastling(move):
		b.castle(move)
	case b.isEnPassant(move):
		b.Coords[to] = b.Coords[from]
		b.Coords[from] = empty
		direction := 8
		if b.SideToMove == BlackToMove {
			direction = -8
		}
		b.Coords[to-Square(direction)] = empty
	case move.Promotion() != 0:
		promoteTo := move.Promotion()
		offset := uint8(0)
		if b.SideToMove != WhiteToMove {
			offset = 6
		}
		switch promoteTo {
		case 'q':
			promoteTo = Q + offset
		case 'n':
			promoteTo = N + offset
		case 'r':
			promoteTo = R + offset
		case 'b':
			promoteTo = B + offset
		}

		b.ZobristPromotion(move)
		b.Coords[to] = promoteTo
		b.Coords[from] = empty
	default:
		b.ZobristSimpleMove(move)
		b.move(move)
	}

	b.updateEnPassantTarget(move)
	b.updateCastlingRights(move)
	b.updateSideToMove()
}

func (b *Board) move(move Move) {
	b.Coords[move.To()] = b.Coords[move.From()]
	b.Coords[move.From()] = empty
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

func (b *Board) castle(move Move) {
	switch move {
	case WCastleKing:
		b.move(WCastleKing)
		b.move(WCastleKingRook)
	case WCastleQueen:
		b.move(WCastleQueen)
		b.move(WCastleQueenRook)
	case BCastleKing:
		b.move(BCastleKing)
		b.move(BCastleKingRook)
	case BCastleQueen:
		b.move(BCastleQueen)
		b.move(BCastleQueenRook)
	}
}

func CoordInBounds(c Square) bool {
	return c >= 0 && c < 64
}

func fileToCoord(file rune) int {
	for i, f := range Files {
		if f == string(file) {
			return i
		}
	}
	return 0
}

func (b *Board) PlayMoves(moves string) {
	moveSlice := strings.Fields(moves)
	for _, move := range moveSlice {
		b.MoveLongAlg(MoveFromString(move))
	}
}
