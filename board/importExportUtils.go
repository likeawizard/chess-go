package board

import (
	"fmt"
	"strconv"
)

func (b Board) ExportFEN() string {
	var fen string
	var emptySquaresCounter, piece int
	for r := len(b.coords) - 1; r >= 0; r-- {
		for f := 0; f < len(b.coords); f++ {
			piece = b.coords[f][r]
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
	fen += fmt.Sprintf(" %s %s %s %d %d", b.sideToMove, b.castlingRights, b.enPassantTarget, b.halfMoveCounter, b.fullMoveCounter)
	return fen
}

func (b *Board) ImportFEN(fen string) {
	var (
		f         = 0
		r         = 7
		chars     = []rune(fen)
		delimiter int
	)

	b.coords = [8][8]int{}
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
				b.coords[f][r] = piece
				f++
			}
		} else {
			f += offset
		}
	}
	fen = fen[delimiter:]
	b.sideToMove = fen[:1]
	b.fullMoveCounter, _ = strconv.Atoi(fen[len(fen)-1:])
	b.halfMoveCounter, _ = strconv.Atoi(fen[len(fen)-3 : len(fen)-2])

	fen = fen[2 : len(fen)-4]

	chars = []rune(fen)

	for i, char := range chars {
		if string(char) == " " {
			delimiter = i + 1
		}
	}
	b.castlingRights = fen[:delimiter-1]
	b.enPassantTarget = fen[delimiter:]
}
