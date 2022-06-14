package board

import (
	"fmt"
	"strconv"
	"strings"
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

func (b *Board) ImportFEN(fen string) error {
	fields := strings.Fields(fen)
	if len(fields) != 6 {
		return fmt.Errorf("FEN must contain six fields - '%s'", fen)
	}
	position := fields[0]
	sideToMove, castling, enPassant, halfMove, fullMove := fields[1], fields[2], fields[3], fields[4], fields[5]

	var err error

	b.Coords, err = parsePosition(position)
	if err != nil {
		return err
	}

	b.SideToMove = sideToMove
	b.FullMoveCounter, err = strconv.Atoi(fullMove)
	if err != nil {
		return err
	}

	b.HalfMoveCounter, err = strconv.Atoi(halfMove)
	if err != nil {
		return err
	}

	b.CastlingRights = castling
	b.EnPassantTarget = enPassant

	return nil
}

func parsePosition(position string) ([8][8]int, error) {
	var (
		f = 0
		r = 7
	)
	c := [8][8]int{}
	for _, char := range position {
		symbol := string(char)
		offset, err := strconv.Atoi(symbol)
		if err != nil {
			if char == '/' {
				f = 0
				r--
			} else if char == ' ' {
				break
			} else {
				piece := PieceSymbolToInt(symbol)
				c[f][r] = piece
				f++
			}
		} else {
			f += offset
		}
	}
	return c, nil
}
