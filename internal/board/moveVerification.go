package board

import (
	"strings"
)

func (b Board) VerifyMove(longalg string) bool {
	from, _ := longAlgToCoords(longalg)
	moves, captures := b.GetAvailableMoves(from)

	_, color := GetPiece(b, from)

	if color != b.SideToMove {
		return false
	}

	for _, move := range moves {
		if move == longalg {
			return true
		}
	}

	for _, capture := range captures {
		if capture == longalg {
			return true
		}
	}

	return false

}

func isOpponentPiece(b Board, source, target Coord) bool {
	piece, color := GetPiece(b, target)
	_, ownColor := GetPiece(b, source)
	return piece != 0 && color != ownColor
}

func (b Board) GetAvailableMoves(c Coord) (availableMoves, availableCaptures []string) {
	return b.GetAvailableMovesRaw(c, false)
}

func (b Board) GetAvailableMovesExcludeCastling(c Coord) (availableMoves, availableCaptures []string) {
	return b.GetAvailableMovesRaw(c, true)
}

func (b Board) GetAvailableMovesRaw(c Coord, excludeCastling bool) (availableMoves, availableCaptures []string) {
	piece := b.AccessCoord(c)
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

func (b Board) GetPawnMoves(c Coord) (moves, captures []string) {
	isWhite := b.Coords[c.File][c.Rank] <= PieceOffset
	pawnCaptures := [2][2]int{{1, 1}, {-1, 1}}
	var isFirstMove bool
	var direction = 1

	if isWhite {
		isFirstMove = c.Rank == 1
	} else {
		isFirstMove = c.Rank == 6
		direction *= -1
	}

	if b.AccessCoord(Coord{c.File, c.Rank + direction}) == 0 {
		moves = append(moves, CoordsToMove(c, Coord{c.File, c.Rank + direction}))
	}

	if isFirstMove && b.AccessCoord(Coord{c.File, c.Rank + direction}) == 0 && b.AccessCoord(Coord{c.File, c.Rank + 2*direction}) == 0 {
		moves = append(moves, CoordsToMove(c, Coord{c.File, c.Rank + 2*direction}))
	}

	for i := 0; i < 2; i++ {
		targetCoord := Coord{c.File + pawnCaptures[i][0], c.Rank + (direction * pawnCaptures[i][1])}
		if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
			captures = append(captures, CoordsToMove(c, targetCoord))
		}
	}

	for i := 0; i < 2; i++ {
		targetCoord := Coord{c.File + pawnCaptures[i][0], c.Rank + (direction * pawnCaptures[i][1])}
		if CoordInBounds(targetCoord) && CoordToAlg(targetCoord) == b.EnPassantTarget {
			captures = append(captures, CoordsToMove(c, targetCoord))
		}
	}

	return
}

func (b Board) GetBishopMoves(c Coord) (moves, captures []string) {
	var ul, ur, dl, dr bool
	var targetCoord Coord
	for i := 1; i < 7; i++ {
		if !ul {
			targetCoord = Coord{c.File + i, c.Rank + i}
			if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
				ul = true
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				ul = true
			}
		}
		if !ur {
			targetCoord = Coord{c.File + i, c.Rank - i}
			if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
				ur = true
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				ur = true
			}
		}

		if !dl {
			targetCoord = Coord{c.File - i, c.Rank + i}
			if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
				dl = true
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				dl = true
			}
		}

		if !dr {
			targetCoord = Coord{c.File - i, c.Rank - i}
			if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
				dr = true
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				dr = true
			}
		}

	}
	return
}

func (b Board) GetKnightMoves(c Coord) (moves, captures []string) {
	var knightMoves = [8][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}}
	var targetCoord Coord

	for i := 0; i < 8; i++ {
		targetCoord = Coord{
			c.File + knightMoves[i][0], c.Rank + knightMoves[i][1],
		}
		if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
			moves = append(moves, CoordsToMove(c, targetCoord))
		} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
			captures = append(captures, CoordsToMove(c, targetCoord))
		}

	}

	return
}

func (b Board) GetRookMoves(c Coord) (moves, captures []string) {
	var u, d, l, r bool
	var targetCoord Coord
	for i := 1; i < 7; i++ {
		if !u {
			targetCoord = Coord{c.File, c.Rank + i}
			if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
				u = true
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				u = true
			}
		}
		if !d {
			targetCoord = Coord{c.File, c.Rank - i}
			if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
				d = true
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				d = true
			}
		}

		if !l {
			targetCoord = Coord{c.File - i, c.Rank}
			if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
				l = true
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				l = true
			}
		}

		if !r {
			targetCoord = Coord{c.File + i, c.Rank}
			if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
				r = true
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				r = true
			}
		}

	}
	return
}

func (b Board) GetQueenMoves(c Coord) (moves, captures []string) {
	bishopMoves, bishopCaptures := b.GetBishopMoves(c)
	rookMoves, rookCaptures := b.GetRookMoves(c)
	moves = append(bishopMoves, rookMoves...)
	captures = append(bishopCaptures, rookCaptures...)
	return
}

func (b Board) GetKingMoves(c Coord, excludeCastling bool) (moves, captures []string) {
	var kingMoves = [8][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	var targetCoord Coord

	for i := 0; i < 8; i++ {
		targetCoord = Coord{
			c.File + kingMoves[i][0], c.Rank + kingMoves[i][1],
		}

		if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == 0 {
			moves = append(moves, CoordsToMove(c, targetCoord))
		} else if CoordInBounds(targetCoord) && isOpponentPiece(b, c, targetCoord) {
			captures = append(captures, CoordsToMove(c, targetCoord))
		}
	}

	if excludeCastling {
		return
	}

	if b.SideToMove == WhiteToMove {
		m, c := b.GetMovesNoCastling(BlackToMove)
		moveDest := movesToDestinationSquaresString(m)
		captureDest := movesToDestinationSquaresString(c)
		var (
			c1 = Coord{2, 0}
			d1 = Coord{3, 0}
			f1 = Coord{5, 0}
			g1 = Coord{6, 0}
		)

		if strings.Contains(b.CastlingRights, wOOO) && b.AccessCoord(c1) == 0 && b.AccessCoord(d1) == 0 &&
			!strings.Contains(moveDest, "c1") && !strings.Contains(moveDest, "d1") && !strings.Contains(captureDest, "e1") {
			moves = append(moves, "e1c1")
		}

		if strings.Contains(b.CastlingRights, wOO) && b.AccessCoord(f1) == 0 && b.AccessCoord(g1) == 0 &&
			!strings.Contains(moveDest, "f1") && !strings.Contains(moveDest, "g1") && !strings.Contains(captureDest, "e1") {
			moves = append(moves, "e1g1")
		}
	} else {
		m, c := b.GetMovesNoCastling(WhiteToMove)
		moveDest := movesToDestinationSquaresString(m)
		captureDest := movesToDestinationSquaresString(c)
		var (
			c8 = Coord{2, 7}
			d8 = Coord{3, 7}
			f8 = Coord{5, 7}
			g8 = Coord{6, 7}
		)

		if strings.Contains(b.CastlingRights, bOOO) && b.AccessCoord(c8) == 0 && b.AccessCoord(d8) == 0 &&
			!strings.Contains(moveDest, "c8") && !strings.Contains(moveDest, "d8") && !strings.Contains(captureDest, "e8") {
			moves = append(moves, "e8c8")
		}

		if strings.Contains(b.CastlingRights, bOO) && b.AccessCoord(f8) == 0 && b.AccessCoord(g8) == 0 &&
			!strings.Contains(moveDest, "f8") && !strings.Contains(moveDest, "g8") && !strings.Contains(captureDest, "e8") {
			moves = append(moves, "e8g8")
		}
	}
	return
}

func movesToDestinationSquaresString(moves []string) (destination string) {
	for _, move := range moves {
		destination += move[2:]
	}
	return destination
}

func (b Board) isCastling(move string) bool {
	for _, castlingMove := range CastlingMoves {
		if move == castlingMove {
			return true
		}
	}
	return false
}

func (b Board) isEnPassant(move string) bool {
	from, to := move[:2], move[2:]
	piece := b.AccessCoord(AlgToCoord(from))
	return (piece == 1 || piece == 7) && to == b.EnPassantTarget
}
