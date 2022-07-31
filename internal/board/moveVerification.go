package board

// N, S, W, E, NW, NE, SW, SE
// 0:4 ranks & files, 4:8 diagonals
var compass = []Square{8, -8, -1, 1, 7, 9, -9, -7}

// Number of squares to the edge in compass direction
var compassBlock = [][]Square{}

var knightMoves = [][]Square{}

//pre calculate distances in all compass directions and possible knight jumps for every square
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

func (b *Board) isOpponentPiece(isWhite bool, sq Square) bool {
	theirPiece := b.Coords[sq]
	if theirPiece == 0 {
		return false
	} else if isWhite && theirPiece > 6 || !isWhite && theirPiece <= 6 {
		return true
	} else {
		return false
	}
}

func (b *Board) IsInCheckAfterMove(move Move) bool {
	umake := b.MoveLongAlg(move)
	defer umake()
	return b.IsInCheck(!b.IsWhite)
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

func (b *Board) GetMovesForPiece(c Square, pin Move) (availableMoves, availableCaptures []Move) {
	piece := b.Coords[c] % PieceOffset

	switch piece {
	case P:
		return b.GetPawnMoves(c)
	case B:
		return b.GetBishopMoves(c, pin)
	case N:
		// pinned knights can't unpin themselves
		if pin != 0 {
			return
		}
		return b.GetKnightMoves(c)
	case R:
		return b.GetRookMoves(c, pin)
	case Q:
		return b.GetQueenMoves(c, pin)
	case 0:
		// King%6 == 0. Will fail if square is empty
		return b.GetKingMoves(c)
	default:
		return
	}
}

func (b *Board) IsInCheck(isWhite bool) bool {
	king := b.GetKing(isWhite)
	return b.IsThretened(isWhite, king)
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
		if CoordInBounds(target) && b.isOpponentPiece(b.IsWhite, target) {
			captures = append(captures, MoveFromSquares(c, target))
			if target < 15 || target > 55 {
				hasPromotion = true
			}
		}
		target = c + direction - 1
		if CoordInBounds(target) && b.isOpponentPiece(b.IsWhite, target) {
			captures = append(captures, MoveFromSquares(c, target))
			if target < 15 || target > 55 {
				hasPromotion = true
			}
		}
	}

	if c%8 == 0 {
		target = c + direction + 1
		if CoordInBounds(target) && b.isOpponentPiece(b.IsWhite, target) {
			captures = append(captures, MoveFromSquares(c, target))
			if target < 15 || target > 55 {
				hasPromotion = true
			}
		}
	}

	if c%8 == 7 {
		target = c + direction - 1
		if CoordInBounds(target) && b.isOpponentPiece(b.IsWhite, target) {
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

// Returns compass directions allowed by pin
func GetCompassPinned(pin Move) []bool {
	// No pin all directions allowed
	if pin == 0 {
		return []bool{true, true, true, true, true, true, true, true}
	}

	// Pinnded along N & S
	if (pin.From()-pin.To())%8 == 0 {
		return []bool{true, true, false, false, false, false, false, false}
	}

	// Pinnded along W & E
	if pin.From()/8 == pin.To()/8 {
		return []bool{false, false, true, true, false, false, false, false}
	}

	// Pinned along NW & SE
	if pin.From()%7 == pin.To()%7 {
		return []bool{false, false, false, false, true, false, false, true}
	}

	return []bool{false, false, false, false, false, true, true, false}

}

func (b *Board) GetSlidingMoves(c Square, mode SlideMode, pin Move) (moves, captures []Move) {
	var target Square

	compassMin, compassMax := 0, 8
	switch mode {
	case BISHOP:
		compassMin = 4
	case ROOK:
		compassMax = 4
	}

	pinnedDirections := GetCompassPinned(pin)
	for dirIdx := compassMin; dirIdx < compassMax; dirIdx++ {
		if !pinnedDirections[dirIdx] {
			continue
		}
		for i := Square(1); i <= compassBlock[c][dirIdx]; i++ {
			target = c + i*compass[dirIdx]
			if b.Coords[target] == 0 {
				moves = append(moves, MoveFromSquares(c, target))
			} else if b.isOpponentPiece(b.IsWhite, target) {
				captures = append(captures, MoveFromSquares(c, target))
				break
			} else {
				break
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

func (b *Board) GetBishopMoves(c Square, pin Move) (moves, captures []Move) {
	return b.GetSlidingMoves(c, BISHOP, pin)
}

func (b *Board) GetRookMoves(c Square, pin Move) (moves, captures []Move) {
	return b.GetSlidingMoves(c, ROOK, pin)
}

func (b *Board) GetQueenMoves(c Square, pin Move) (moves, captures []Move) {
	return b.GetSlidingMoves(c, QUEEN, pin)
}

func (b *Board) GetKnightMoves(c Square) (moves, captures []Move) {
	for _, knightMove := range knightMoves[c] {
		if b.Coords[knightMove] == 0 {
			moves = append(moves, MoveFromSquares(c, knightMove))
		} else if b.isOpponentPiece(b.IsWhite, knightMove) {
			captures = append(captures, MoveFromSquares(c, knightMove))
		}
	}
	return
}

func (b *Board) GetKingMoves(c Square) (moves, captures []Move) {
	for i := 0; i < 8; i++ {
		if compassBlock[c][i] == 0 {
			continue
		}

		if b.IsThretened(b.IsWhite, c+compass[i]) {
			continue
		}

		if b.Coords[c+compass[i]] == 0 {
			moves = append(moves, MoveFromSquares(c, c+compass[i]))
		} else if b.isOpponentPiece(b.IsWhite, c+compass[i]) {
			captures = append(captures, MoveFromSquares(c, c+compass[i]))
		}
	}

	if b.IsInCheck(b.IsWhite) {
		return
	}

	bSq := Square(1)
	cSq := Square(2)
	dSq := Square(3)
	fSq := Square(5)
	gSq := Square(6)
	kingSideCastle := WCastleKing
	kingSideCastleRights := WOO
	queenSideCastle := WCastleQueen
	queenSideCastleRights := WOOO

	if !b.IsWhite {
		bSq = Square(57)
		cSq = Square(58)
		dSq = Square(59)
		fSq = Square(61)
		gSq = Square(62)
		kingSideCastle = BCastleKing
		kingSideCastleRights = BOO
		queenSideCastle = BCastleQueen
		queenSideCastleRights = BOOO
	}

	if b.CastlingRights&queenSideCastleRights != 0 && b.Coords[cSq] == 0 && b.Coords[bSq] == 0 && b.Coords[dSq] == 0 &&
		!b.IsThretened(b.IsWhite, cSq) && !b.IsThretened(b.IsWhite, dSq) {
		moves = append(moves, queenSideCastle)
	}

	if b.CastlingRights&kingSideCastleRights != 0 && b.Coords[fSq] == 0 && b.Coords[gSq] == 0 &&
		!b.IsThretened(b.IsWhite, fSq) && !b.IsThretened(b.IsWhite, gSq) {
		moves = append(moves, kingSideCastle)
	}

	return
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

// Calculate absolute pins to the king to determine move legality of pinned piece. Only moves in the pin direction are allowed
// Move represents square combinations of from -pinned piece and to - the pinner.
func (b *Board) GetPins(isWhite bool) []Move {
	king := b.GetKing(isWhite)
	offset := uint8(6)
	if !b.IsWhite {
		offset = 0
	}

	// idx 0..3 rook pins, 4..7 rooks
	doesPin := func(idx int, piece uint8) bool {
		if idx < 4 && (piece == Q+offset || piece == R+offset) {
			return true
		} else if idx >= 4 && (piece == Q+offset || piece == B+offset) {
			return true
		} else {
			return false
		}
	}

	var target Square
	pins := make([]Move, 0)

	for dirIdx := 0; dirIdx < 8; dirIdx++ {
		pinned := Square(-1)
		for i := Square(1); i <= compassBlock[king][dirIdx]; i++ {
			target = king + i*compass[dirIdx]
			if b.Coords[target] == 0 {
				continue
			}

			if !b.isOpponentPiece(isWhite, target) && pinned == -1 {
				// set first friendly piece in direction as pinned
				pinned = target
			} else if !b.isOpponentPiece(isWhite, target) && pinned != -1 {
				// if second piece is also friendly. stop no pins in direction possible
				break
			} else if pinned != 0 && doesPin(dirIdx, b.Coords[target]) {
				// if piece is correct pinner type add pin
				pins = append(pins, MoveFromSquares(pinned, target))
				break
			} else {
				// opponent piece cant pin in this direction
				break
			}
		}
	}

	return pins
}

// Determine if a square is thretened by the opposition
// checks for opposing p
func (b *Board) IsThretened(isWhite bool, sq Square) bool {
	var target Square
	pawnDirection := Square(-8)
	offset := uint8(0)
	if isWhite {
		pawnDirection = 8
		offset = 6
	}

	target = sq + 1 + pawnDirection
	if sq%8 != 7 && CoordInBounds(target) && b.Coords[target] == P+offset {
		return true
	}
	target = sq - 1 + pawnDirection
	if sq%8 != 0 && CoordInBounds(target) && b.Coords[target] == P+offset {
		return true
	}

	for _, knightMove := range knightMoves[sq] {
		if CoordInBounds(knightMove) && b.Coords[knightMove] == N+offset {
			return true
		}
	}

	isThreat := func(idx int, distance Square, piece uint8) bool {
		if idx < 4 && (piece == Q+offset || piece == R+offset || (distance == 1 && piece == K+offset)) {
			return true
		} else if idx >= 4 && (piece == Q+offset || piece == B+offset || (distance == 1 && piece == K+offset)) {
			return true
		} else {
			return false
		}
	}

	for dirIdx := 0; dirIdx < 8; dirIdx++ {
		for i := Square(1); i <= compassBlock[sq][dirIdx]; i++ {
			target = sq + i*compass[dirIdx]
			if b.Coords[target] == 0 {
				continue
			}

			if !b.isOpponentPiece(isWhite, target) || !isThreat(dirIdx, i, b.Coords[target]) {
				break
			} else if b.isOpponentPiece(isWhite, target) && isThreat(dirIdx, i, b.Coords[target]) {
				return true
			}
		}
	}
	return false
}
