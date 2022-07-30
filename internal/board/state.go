package board

func (b *Board) updateCastlingRights(move Move) {
	if b.CastlingRights == 0 {
		return
	}
	from, to := move.FromTo()

	switch {
	case b.CastlingRights&(WOOO|WOO) != 0 && from == WCastleQueen.From():
		if b.CastlingRights&WOO != 0 {
			b.ZobristCastlingRights(WOO)
		}
		if b.CastlingRights&WOOO != 0 {
			b.ZobristCastlingRights(WOOO)
		}

		b.CastlingRights = b.CastlingRights &^ WOOO
		b.CastlingRights = b.CastlingRights &^ WOO

	case b.CastlingRights&(BOOO|BOO) != 0 && from == BCastleQueen.From():
		if b.CastlingRights&BOOO != 0 {
			b.ZobristCastlingRights(BOOO)
		}
		if b.CastlingRights&BOO != 0 {
			b.ZobristCastlingRights(BOO)
		}

		b.CastlingRights = b.CastlingRights &^ BOOO
		b.CastlingRights = b.CastlingRights &^ BOO

	case b.CastlingRights&WOOO != 0 && (from == WCastleQueenRook.From() || to == WCastleQueenRook.From()):
		b.ZobristCastlingRights(WOOO)
		b.CastlingRights = b.CastlingRights &^ WOOO

	case b.CastlingRights&WOO != 0 && (from == WCastleKingRook.From() || to == WCastleKingRook.From()):
		b.ZobristCastlingRights(WOO)
		b.CastlingRights = b.CastlingRights &^ WOO

	case b.CastlingRights&BOOO != 0 && (from == BCastleQueenRook.From() || to == BCastleQueenRook.From()):
		b.ZobristCastlingRights(BOOO)
		b.CastlingRights = b.CastlingRights &^ BOOO

	case b.CastlingRights&BOO != 0 && (from == BCastleKingRook.From() || to == BCastleKingRook.From()):
		b.ZobristCastlingRights(BOO)
		b.CastlingRights = b.CastlingRights &^ BOO
	}
}

func (b *Board) updateEnPassantTarget(move Move) {
	from, to := move.FromTo()
	piece := b.Coords[to]
	isPawnMove := piece == P || piece == p
	isDoubleMove := from-to == 16 || from-to == -16

	if b.EnPassantTarget != -1 {
		b.ZobristEnPassant(b.EnPassantTarget)
	}

	if isDoubleMove && isPawnMove {
		b.EnPassantTarget = Square((from + to) / 2)
		b.ZobristEnPassant(b.EnPassantTarget)
	} else {
		b.EnPassantTarget = -1
	}
}

func (b *Board) updateSideToMove() {
	b.ZobristSideToMove()
	b.IsWhite = !b.IsWhite

	//Increment moves after black moves
	if b.IsWhite {
		b.FullMoveCounter++
	}
}
