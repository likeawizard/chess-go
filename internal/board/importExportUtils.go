package board

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (b *Board) ExportFEN() string {
	var fen string
	var emptySquaresCounter int
	var piece uint8
	for r := 7; r >= 0; r-- {
		for f := 0; f < 8; f++ {
			piece = b.Coords[f+r*8]
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
	castlingRights := ""
	if b.CastlingRights != 0 {
		if b.CastlingRights&WOO != 0 {
			castlingRights += "K"
		}
		if b.CastlingRights&WOOO != 0 {
			castlingRights += "Q"
		}
		if b.CastlingRights&BOO != 0 {
			castlingRights += "k"
		}
		if b.CastlingRights&BOOO != 0 {
			castlingRights += "q"
		}
	} else {
		castlingRights = "-"
	}

	epString := "-"
	if b.EnPassantTarget != -1 {
		epString = b.EnPassantTarget.String()
	}

	sideToMove := WhiteToMove
	if !b.IsWhite {
		sideToMove = BlackToMove
	}

	fen += fmt.Sprintf(" %c %s %s %d %d", sideToMove, castlingRights, epString, b.HalfMoveCounter, b.FullMoveCounter)
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

	b.IsWhite = sideToMove[0] == WhiteToMove
	fm, err := strconv.Atoi(fullMove)
	if err != nil {
		return err
	}
	b.FullMoveCounter = uint8(fm)

	hm, err := strconv.Atoi(halfMove)
	if err != nil {
		return err
	}
	b.HalfMoveCounter = uint8(hm)

	for _, c := range []byte(castling) {
		switch c {
		case 'K':
			b.CastlingRights = b.CastlingRights | WOO
		case 'Q':
			b.CastlingRights = b.CastlingRights | WOOO
		case 'k':
			b.CastlingRights = b.CastlingRights | BOO
		case 'q':
			b.CastlingRights = b.CastlingRights | BOOO
		}
	}

	if enPassant != "-" {
		b.EnPassantTarget = SquareFromString(enPassant)
	} else {
		b.EnPassantTarget = -1
	}

	b.Hash = b.SeedHash()

	return nil
}

func parsePosition(position string) ([64]uint8, error) {
	var (
		f = 0
		r = 7
	)
	c := [64]uint8{}
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
				c[f+r*8] = piece
				f++
			}
		} else {
			f += offset
		}
	}
	return c, nil
}

func (b *Board) WritePGNToFile(data string, path string) {
	os.WriteFile(path, []byte(data), 0644)
}

func (b *Board) GeneratePGN(moves []Move) string {
	pgn := ""
	bb := &Board{}
	bb.InitDefault()
	for n, move := range moves {
		if n%2 == 0 {
			pgn += fmt.Sprintf("%d. ", bb.FullMoveCounter)
		}
		pgn += bb.MoveToPretty(move.String()) + " "
		bb.MoveLongAlg(move)
	}
	return pgn
}

func (b *Board) MoveToPretty(move string) (pretty string) {
	var CastlingMoves = [4]string{"e1g1", "e1c1", "e8g8", "e8c8"}
	from, to := MoveFromString(move).FromTo()
	targetPiece := b.Coords[to]
	piece := b.Coords[from]
	all := b.GetLegalMoves()
	switch {
	case piece == P || piece == p:
		pretty = move[2:]
		if move[:1] != move[2:3] {
			pretty = move[:1] + "x" + pretty
		}
	case (piece == K || piece == k) && move == CastlingMoves[0] || move == CastlingMoves[2]:
		return "O-O"
	case (piece == K || piece == k) && move == CastlingMoves[1] || move == CastlingMoves[3]:
		return "O-O-O"
	default:
		pretty = Pieces[(piece-1)%PieceOffset]
		pretty += b.Disambiguate(move, all)
		if targetPiece > 0 {
			pretty += "x"
		}
		pretty += move[2:]
	}
	if len(move) == 5 {
		pretty += "=" + strings.ToUpper(move[4:])
	}

	return
}

func (b *Board) Disambiguate(move string, moves []Move) string {
	dis := ""
	from, to := MoveFromString(move).FromTo()
	for _, m := range moves {
		f, t := m.FromTo()
		if m.String()[:2] == move[:2] || b.Coords[from] != b.Coords[f] {
			continue
		}
		if (from-f)%8 == 0 && to == t {
			dis += move[1:2]
		}
		if f/8 == from/8 && to == t {
			dis += move[0:1]
		}
	}
	return dis
}
