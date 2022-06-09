package board

import (
	"fmt"
	"strconv"
)

func (b *Board) ExportFEN() string {
	var fen string
	var emptySquaresCounter, piece int
	for r := len(b.Coords) - 1; r >= 0; r-- {
		for f := 0; f < len(b.Coords); f++ {
			piece = b.Coords[f][r]
			if piece == 0 {
				emptySquaresCounter++
			} else {
				if emptySquaresCounter > 0 {
					fen += string(rune(emptySquaresCounter + '0'))
					emptySquaresCounter = 0
				}
				fen += Pieces[piece-1]
			}
		}
		if emptySquaresCounter > 0 {
			fen += string(rune(emptySquaresCounter + '0'))
			emptySquaresCounter = 0
		}
		if r > 0 {
			fen += "/"
		}
	}
	fen += fmt.Sprintf(" %s %s %s %d %d", b.SideToMove, b.CastlingRights, b.EnPassantTarget, b.HalfMoveCounter, b.FullMoveCounter)
	return fen
}

func (b *Board) ImportFEN(fen string) {
	var (
		f         = 0
		r         = 7
		chars     = []rune(fen)
		delimiter int
	)

	b.Coords = [8][8]int{}
	for index, char := range chars {
		symbol := string(char)
		offset, err := strconv.Atoi(symbol)
		if err != nil {
			if symbol == "/" {
				f = 0
				r--
			} else if symbol == " " {
				delimiter = index + 1
				break
			} else {
				piece := PieceSymbolToInt(symbol)
				b.Coords[f][r] = piece
				f++
			}
		} else {
			f += offset
		}
	}
	fen = fen[delimiter:]
	b.SideToMove = fen[:1]
	b.FullMoveCounter, _ = strconv.Atoi(fen[len(fen)-1:])
	b.HalfMoveCounter, _ = strconv.Atoi(fen[len(fen)-3 : len(fen)-2])

	fen = fen[2 : len(fen)-4]

	chars = []rune(fen)

	for i, char := range chars {
		if string(char) == " " {
			delimiter = i + 1
		}
	}
	b.CastlingRights = fen[:delimiter-1]
	b.EnPassantTarget = fen[delimiter:]
}
