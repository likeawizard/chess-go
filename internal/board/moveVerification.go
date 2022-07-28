package board

// N, S, W, E, NW, NE, SW, SE
// 0:4 ranks & files, 4:8 diagonals
var compass = []Square{8, -8, -1, 1, 7, 9, -9, -7}

// Number of squares to the edge in compass direction
var compassBlock = [][]Square{}

var knightMoves = [][]Square{}

func init() {
	compassBlock = make([][]Square, 64)
	min := func(a, b Square) Square {
		if a < b {
			return a
		}
		return b
	}
	for i := Square(0); i < 64; i++ {
		f, r := i%8, i/8
		n := 7 - r
		s := r
		w := f
		e := 7 - f
		compassBlock[i] = []Square{n, s, w, e, min(n, w), min(n, e), min(s, w), min(s, e)}
	}
	preCalculateKnightMoves()
}

func preCalculateKnightMoves() {
	// var knightVector = [8]Square{15, 17, 10, 6, -15, -17, -10, -6}
	knightMoves = make([][]Square, 64)
	for c, comp := range compassBlock {
		moves := make([]Square, 0)
		//2NW
		if comp[0] > 1 && comp[2] > 0 {
			moves = append(moves, Square(c)+15)
		}
		//2NE
		if comp[0] > 1 && comp[3] > 0 {
			moves = append(moves, Square(c)+17)
		}
		//2SW
		if comp[1] > 1 && comp[2] > 0 {
			moves = append(moves, Square(c)-17)
		}
		//2SE
		if comp[1] > 1 && comp[3] > 0 {
			moves = append(moves, Square(c)-15)
		}
		//2WN
		if comp[2] > 1 && comp[0] > 0 {
			moves = append(moves, Square(c)+6)
		}
		//2WS
		if comp[2] > 1 && comp[1] > 0 {
			moves = append(moves, Square(c)-10)
		}
		//2EN
		if comp[3] > 1 && comp[0] > 0 {
			moves = append(moves, Square(c)+10)
		}
		//2ES
		if comp[3] > 1 && comp[1] > 0 {
			moves = append(moves, Square(c)-6)
		}
		knightMoves[c] = moves
	}
}

func (b *Board) isOpponentPiece(us, them Square) bool {
	ourPiece, theirPiece := b.Coords[us], b.Coords[them]
	return (ourPiece < 7 && theirPiece >= 7) || (ourPiece >= 7 && theirPiece < 7 && theirPiece != 0)
}

func (b *Board) IsInCheckAfterMove(move Move) bool {
	bb := b.Copy()
	bb.MoveLongAlg(move)
	return bb.IsInCheck(b.IsWhite)
}

func (b *Board) PruneIllegal(moves, captures []Move) ([]Move, []Move) {
	legalMoves := make([]Move, 0)
	legalCaptures := make([]Move, 0)
	for _, move := range moves {
		if !b.IsInCheckAfterMove(move) {
			legalMoves = append(legalMoves, move)
		}
	}
	for _, capture := range captures {
		if !b.IsInCheckAfterMove(capture) {
			legalCaptures = append(legalCaptures, capture)
		}
	}

	return legalMoves, legalCaptures
}

func (b *Board) GetAvailableMoves(c Square) (availableMoves, availableCaptures []Move) {
	return b.GetAvailableMovesRaw(c, false)
}

func (b *Board) GetAvailableMovesExcludeCastling(c Square) (availableMoves, availableCaptures []Move) {
	return b.GetAvailableMovesRaw(c, true)
}

func (b *Board) GetAvailableMovesRaw(c Square, excludeCastling bool) (availableMoves, availableCaptures []Move) {
	piece := b.Coords[c]
	plainPiece := piece % PieceOffset

	switch {
	case plainPiece == 1:
		return b.GetPawnMoves(c)
	case plainPiece == 2:
		return b.GetBishopMoves(c)
	case plainPiece == 3:
		return b.GetKnightMoves(c)
	case plainPiece == 4:
		return b.GetRookMoves(c)
	case plainPiece == 5:
		return b.GetQueenMoves(c)
	case piece == 6 || piece == 12:
		return b.GetKingMoves(c, excludeCastling)
	default:
		return
	}
}

func (b *Board) IsInCheck(isWhite bool) bool {
	var target Square
	pawnDirection := Square(-8)
	offset := uint8(0)
	king := b.GetKing(isWhite)
	if isWhite {
		pawnDirection = 8
		offset = 6
	}

	target = king + 1 + pawnDirection
	if CoordInBounds(target) && b.Coords[target] == P+offset {
		return true
	}
	target = king - 1 + pawnDirection
	if CoordInBounds(target) && b.Coords[target] == P+offset {
		return true
	}

	for _, knightMove := range knightMoves[king] {
		if CoordInBounds(knightMove) && b.Coords[knightMove] == N+offset {
			return true
		}
	}

	_, kingThreats := b.GetKingMoves(king, true)
	_, bishopThreats := b.GetBishopMoves(king)
	_, rookThreats := b.GetRookMoves(king)

	for _, diagAttacks := range kingThreats {
		attacker := b.Coords[diagAttacks.To()]
		if attacker == Q+offset || attacker == K+offset {
			return true
		}
	}

	for _, diagAttacks := range bishopThreats {
		attacker := b.Coords[diagAttacks.To()]
		if attacker == Q+offset || attacker == B+offset {
			return true
		}
	}

	for _, rookAttacker := range rookThreats {
		attacker := b.Coords[rookAttacker.To()]
		if attacker == Q+offset || attacker == R+offset {
			return true
		}
	}
	return false
}

func (b *Board) GetPawnMoves(c Square) (moves, captures []Move) {
	var target Square
	isWhite := b.Coords[c] <= PieceOffset
	var isFirstMove bool
	var direction = Square(8)
	hasPromotion := false

	if isWhite {
		isFirstMove = c >= 8 && c < 16
	} else {
		isFirstMove = c >= 48 && c < 56
		direction = -8
	}

	if CoordInBounds(c+direction) && b.Coords[c+direction] == 0 {
		moves = append(moves, MoveFromSquares(c, c+direction))
		if c+direction < 15 || c+direction > 55 {
			hasPromotion = true
		}
	}

	if isFirstMove && b.Coords[c+direction] == 0 && b.Coords[c+2*direction] == 0 {
		moves = append(moves, MoveFromSquares(c, c+2*direction))
	}

	if c%8 > 0 && c%8 < 7 {
		target = c + direction + 1
		if CoordInBounds(target) && b.isOpponentPiece(c, target) {
			captures = append(captures, MoveFromSquares(c, target))
			if target < 15 || target > 55 {
				hasPromotion = true
			}
		}
		target = c + direction - 1
		if CoordInBounds(target) && b.isOpponentPiece(c, target) {
			captures = append(captures, MoveFromSquares(c, target))
			if target < 15 || target > 55 {
				hasPromotion = true
			}
		}
	}

	if c%8 == 0 {
		target = c + direction + 1
		if CoordInBounds(target) && b.isOpponentPiece(c, target) {
			captures = append(captures, MoveFromSquares(c, target))
			if target < 15 || target > 55 {
				hasPromotion = true
			}
		}
	}

	if c%8 == 7 {
		target = c + direction - 1
		if CoordInBounds(target) && b.isOpponentPiece(c, target) {
			captures = append(captures, MoveFromSquares(c, target))
			if target < 15 || target > 55 {
				hasPromotion = true
			}
		}
	}

	if (c/8 == 3 || c/8 == 4) && b.EnPassantTarget != -1 {
		if c%8 > 0 && c%8 < 7 {
			target = c + direction + 1
			if target == b.EnPassantTarget {
				captures = append(captures, MoveFromSquares(c, target))
			}
			target = c + direction - 1
			if target == b.EnPassantTarget {
				captures = append(captures, MoveFromSquares(c, target))
			}
		}

		if c%8 == 0 {
			target = c + direction + 1
			if target == b.EnPassantTarget {
				captures = append(captures, MoveFromSquares(c, target))
			}
		}

		if c%8 == 7 {
			target = c + direction - 1
			if target == b.EnPassantTarget {
				captures = append(captures, MoveFromSquares(c, target))
			}
		}
	}

	if hasPromotion {
		moves, captures = b.addPawnPromotion(moves, captures)
	}

	return
}

func (b *Board) addPawnPromotion(moves, captures []Move) ([]Move, []Move) {
	processMoves := func(moves []Move) []Move {
		var m []Move
		for _, move := range moves {
			to := move.To()
			if to/8 == 7 || to/8 == 0 {
				m = append(m, move.SetPromotion('q'), move.SetPromotion('r'), move.SetPromotion('n'), move.SetPromotion('b'))
			} else {
				m = append(m, move)
			}

		}
		return m
	}
	return processMoves(moves), processMoves(captures)
}

func (b *Board) GetSlidingMoves(c Square, mode SlideMode) (moves, captures []Move) {
	var target Square

	compassMin, compassMax := 0, 8
	switch mode {
	case BISHOP:
		compassMin = 4
	case ROOK:
		compassMax = 4
	}

	for dirIdx := compassMin; dirIdx < compassMax; dirIdx++ {
		for i := Square(1); i <= compassBlock[c][dirIdx]; i++ {
			target = c + i*compass[dirIdx]
			if CoordInBounds(target) {
				if b.Coords[target] == 0 {
					moves = append(moves, MoveFromSquares(c, target))
				} else if b.isOpponentPiece(c, target) {
					captures = append(captures, MoveFromSquares(c, target))
					break
				} else {
					break
				}
			}
		}
	}

	return
}

type SlideMode byte

const (
	ROOK SlideMode = iota
	BISHOP
	QUEEN
)

func (b *Board) GetBishopMoves(c Square) (moves, captures []Move) {
	return b.GetSlidingMoves(c, BISHOP)
}

func (b *Board) GetRookMoves(c Square) (moves, captures []Move) {
	return b.GetSlidingMoves(c, ROOK)
}

func (b *Board) GetQueenMoves(c Square) (moves, captures []Move) {
	return b.GetSlidingMoves(c, QUEEN)
}

func (b *Board) GetKnightMoves(c Square) (moves, captures []Move) {
	for _, knightMove := range knightMoves[c] {
		if b.Coords[knightMove] == 0 {
			moves = append(moves, MoveFromSquares(c, knightMove))
		} else if b.isOpponentPiece(c, knightMove) {
			captures = append(captures, MoveFromSquares(c, knightMove))
		}
	}
	return
}

func (b *Board) GetKingMoves(c Square, excludeCastling bool) (moves, captures []Move) {
	for i := 0; i < 8; i++ {
		if compassBlock[c][i] == 0 {
			continue
		}
		if CoordInBounds(c+compass[i]) && b.Coords[c+compass[i]] == 0 {
			moves = append(moves, MoveFromSquares(c, c+compass[i]))
		} else if CoordInBounds(c+compass[i]) && b.isOpponentPiece(c, c+compass[i]) {
			captures = append(captures, MoveFromSquares(c, c+compass[i]))
		}
	}

	if excludeCastling {
		return
	}

	if b.IsWhite {
		m, c := b.GetMovesNoCastling(false)
		moveDest := movesToDestinationSquaresString(m)
		captureDest := movesToDestinationSquaresString(c)
		var (
			b1 = Square(1)
			c1 = Square(2)
			d1 = Square(3)
			e1 = Square(4)
			f1 = Square(5)
			g1 = Square(6)

			b2 = Square(9)
			c2 = Square(10)
			e2 = Square(12)
			g2 = Square(14)
			h2 = Square(15)
		)

		if b.CastlingRights&WOOO != 0 && b.Coords[c1] == 0 && b.Coords[b1] == 0 && b.Coords[d1] == 0 &&
			b.Coords[b2] != p && b.Coords[c2] != p && b.Coords[e2] != p &&
			!containsSqaure(moveDest, c1) && !containsSqaure(moveDest, d1) && !containsSqaure(captureDest, e1) {
			moves = append(moves, WCastleQueen)
		}

		if b.CastlingRights&WOO != 0 && b.Coords[f1] == 0 && b.Coords[g1] == 0 &&
			b.Coords[g2] != p && b.Coords[h2] != p && b.Coords[e2] != p &&
			!containsSqaure(moveDest, f1) && !containsSqaure(moveDest, g1) && !containsSqaure(captureDest, e1) {
			moves = append(moves, WCastleKing)
		}
	} else {
		m, c := b.GetMovesNoCastling(true)
		moveDest := movesToDestinationSquaresString(m)
		captureDest := movesToDestinationSquaresString(c)
		var (
			b8 = Square(57)
			c8 = Square(58)
			d8 = Square(59)
			e8 = Square(60)
			f8 = Square(61)
			g8 = Square(62)

			b7 = Square(49)
			c7 = Square(50)
			e7 = Square(52)
			g7 = Square(54)
			h7 = Square(55)
		)

		if b.CastlingRights&BOOO != 0 && b.Coords[c8] == 0 && b.Coords[b8] == 0 && b.Coords[d8] == 0 &&
			b.Coords[b7] != P && b.Coords[c7] != P && b.Coords[e7] != P &&

			!containsSqaure(moveDest, c8) && !containsSqaure(moveDest, d8) && !containsSqaure(captureDest, e8) {
			moves = append(moves, BCastleQueen)
		}

		if b.CastlingRights&BOO != 0 && b.Coords[f8] == 0 && b.Coords[g8] == 0 &&
			b.Coords[g7] != P && b.Coords[h7] != P && b.Coords[e7] != P &&
			!containsSqaure(moveDest, f8) && !containsSqaure(moveDest, g8) && !containsSqaure(captureDest, e8) {
			moves = append(moves, BCastleKing)
		}
	}
	return
}

func containsSqaure(squares []Square, needle Square) bool {
	for i := 0; i < len(squares); i++ {
		if needle == squares[i] {
			return true
		}
	}
	return false
}

func movesToDestinationSquaresString(moves []Move) (destination []Square) {
	for _, move := range moves {
		destination = append(destination, move.To())
	}
	return destination
}

func (b *Board) IsCastling(move Move) bool {
	if b.CastlingRights&CASTLING_ALL == 0 {
		return false
	}
	from := move.From()
	king := b.Coords[from]
	if king != K && king != k {
		return false
	}

	if move == WCastleQueen && b.CastlingRights&WOOO != 0 {
		return true
	}

	if move == WCastleKing && b.CastlingRights&WOO != 0 {
		return true
	}

	if move == BCastleQueen && b.CastlingRights&BOOO != 0 {
		return true
	}

	if move == BCastleKing && b.CastlingRights&BOO != 0 {
		return true
	}
	return false
}

func (b *Board) isEnPassant(move Move) bool {
	from, to := move.FromTo()
	piece := b.Coords[from]
	return (piece == 1 || piece == 7) && to == b.EnPassantTarget
}
