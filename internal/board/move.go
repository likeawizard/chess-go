package board

import (
	"fmt"
	"strconv"
)

var (
	// Castling moves. Used for recognizing castling and moving king during castling
	WCastleKing  = MoveFromString("e1g1")
	WCastleQueen = MoveFromString("e1c1")
	BCastleKing  = MoveFromString("e8g8")
	BCastleQueen = MoveFromString("e8c8")

	// Complimentary castling moves. Used during castling to reposition rook
	WCastleKingRook  = MoveFromString("h1f1")
	WCastleQueenRook = MoveFromString("a1d1")
	BCastleKingRook  = MoveFromString("h8f8")
	BCastleQueenRook = MoveFromString("a8d8")
)

// 0..7 a1 to h1
// 0..63 to a1 to h8 mapping
type Square int

func SquareFromString(s string) Square {
	file := int(s[0] - 'a')
	rank, _ := strconv.Atoi(s[1:])
	rank--
	return Square(file + rank*8)
}

func (s Square) String() string {
	rank := s/8 + 1
	file := s % 8
	return fmt.Sprintf("%c%d", file+'a', rank)
}

// MSB ----------------- LSB
//     prom -- from -- to
// prom - byte, from and to 6bit = 20bit
type Move int

func MoveFromString(s string) Move {
	from := SquareFromString(s[:2]) << 6
	to := SquareFromString(s[2:4])
	promotion := 0
	if len(s) == 5 {
		promotion = int(s[4]) << 12
	}
	return Move(from + to + Square(promotion))
}

func MoveFromSquares(from, to Square) Move {
	from = from << 6
	return Move(from + to)
}

func (m Move) From() Square {
	return Square(m>>6) & 63
}

func (m Move) To() Square {
	return Square(m) & 63
}

func (m Move) Promotion() uint8 {
	return uint8(m >> 12)
}

func (m Move) SetPromotion(prom uint8) Move {
	m &= 4095
	m += Move(prom) << 12
	return m
}

func (m Move) FromTo() (Square, Square) {
	return Square(m>>6) & 63, Square(m) & 63
}

func (m Move) String() string {
	if m.Promotion() != 0 {
		return fmt.Sprintf("%v%v%c", Square(m.From()), Square(m.To()), m.Promotion())
	} else {
		return fmt.Sprintf("%v%v", Square(m.From()), Square(m.To()))
	}
}
