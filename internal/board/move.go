package board

import (
	"fmt"
	"strconv"
)

var (
	WCastleKing  = MoveFromString("e1g1")
	WCastleQueen = MoveFromString("e1c1")
	BCastleKing  = MoveFromString("e8g8")
	BCastleQueen = MoveFromString("e8c8")
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

func SquareFromCoord(c Coord) Square {
	return Square(c.File + c.Rank*8)
}

func (s Square) String() string {
	rank := s/8 + 1
	file := s % 8
	return fmt.Sprintf("%c%d", file+'a', rank)
}

func (s Square) ToCoord() Coord {
	return Coord{File: int(s % 8), Rank: int(s / 8)}
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
	m &= (63<<6 + 63)
	m += (Move(prom) << 12)
	return m
}

func (m Move) ToCoords() (Coord, Coord) {
	return m.From().ToCoord(), m.To().ToCoord()
}

func CoordsToMove(from, to Coord) Move {
	ff := SquareFromCoord(from) << 6
	tt := SquareFromCoord(to)
	return Move(ff + tt)
}

func (m Move) String() string {
	return fmt.Sprintf("%v%v%c", Square(m.From()), Square(m.To()), m.Promotion())
}
