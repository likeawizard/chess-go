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
		IsWhite:         b.IsWhite,
		CastlingRights:  b.CastlingRights,
		EnPassantTarget: b.EnPassantTarget,
		HalfMoveCounter: b.HalfMoveCounter,
		FullMoveCounter: b.HalfMoveCounter,
	}
}

type UnMakeMove func()
type UnMakeMoveOptions struct {
	isWhite         bool
	isCastling      bool
	isEnPassant     bool
	enPassantTarget Square
	cRights         CastlingRights
	targetPiece     uint8
	isPromotion     bool
}

func (b *Board) getUnmake(move Move, opts UnMakeMoveOptions) UnMakeMove {
	unmove := move.Reverse()
	var umake UnMakeMove = func() {
		direction := Square(8)
		if opts.isWhite {
			direction = -8
		}
		b.move(unmove)
		if opts.isEnPassant {
			b.Coords[opts.enPassantTarget+direction] = opts.targetPiece
		} else {
			b.Coords[move.To()] = opts.targetPiece
		}

		if opts.isPromotion {
			if opts.isWhite {
				b.Coords[move.From()] = P
			} else {
				b.Coords[move.From()] = p
			}
		}

		if opts.isCastling {
			switch move {
			case WCastleKing:
				b.move(WCastleKingRook.Reverse())
			case WCastleQueen:
				b.move(WCastleQueenRook.Reverse())
			case BCastleKing:
				b.move(BCastleKingRook.Reverse())
			case BCastleQueen:
				b.move(BCastleQueenRook.Reverse())
			}

		}
		if !opts.isWhite {
			b.FullMoveCounter--
		}

		b.EnPassantTarget = opts.enPassantTarget
		b.CastlingRights = opts.cRights
		b.IsWhite = !b.IsWhite
	}

	return umake
}

func (b *Board) MoveLongAlg(move Move) UnMakeMove {
	from, to := move.FromTo()
	unmake := UnMakeMoveOptions{
		isWhite:         b.IsWhite,
		enPassantTarget: b.EnPassantTarget,
		cRights:         b.CastlingRights,
		targetPiece:     b.Coords[to],
	}
	switch {
	case b.IsCastling(move):
		unmake.isCastling = true
		b.castle(move)
	case b.isEnPassant(move):
		unmake.isEnPassant = true
		b.Coords[to] = b.Coords[from]
		b.Coords[from] = empty
		direction := Square(8)
		if !b.IsWhite {
			direction = -8
		}
		unmake.targetPiece = b.Coords[to-direction]
		b.Coords[to-direction] = empty
	case move.Promotion() != 0:
		unmake.isPromotion = true
		promoteTo := move.Promotion()
		offset := uint8(0)
		if !b.IsWhite {
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

	return b.getUnmake(move, unmake)
}

func (b *Board) move(move Move) {
	b.Coords[move.To()] = b.Coords[move.From()]
	b.Coords[move.From()] = empty
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

func (b *Board) PlayMoves(moves string) {
	moveSlice := strings.Fields(moves)
	for _, move := range moveSlice {
		b.MoveLongAlg(MoveFromString(move))
	}
}
