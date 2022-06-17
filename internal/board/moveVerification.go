package board

import (
	"fmt"
	"strings"
)

func (b *Board) VerifyMove(longalg string) bool {
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

func isOpponentPiece(b *Board, source, target Coord) bool {
	piece, color := GetPiece(b, target)
	_, ownColor := GetPiece(b, source)
	return piece != 0 && color != ownColor
}

func (b *Board) IsInCheckAfterMove(move string) bool {
	bb := b.Copy()
	side := bb.SideToMove
	bb.MoveLongAlg(move)
	return bb.IsInCheck(side)
}

func (b *Board) PruneIllegal(moves, captures []string) ([]string, []string) {
	legalMoves := make([]string, 0)
	legalCaptures := make([]string, 0)
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

func (b *Board) GetAvailableMoves(c Coord) (availableMoves, availableCaptures []string) {
	return b.GetAvailableMovesRaw(c, false)
}

func (b *Board) GetAvailableMovesExcludeCastling(c Coord) (availableMoves, availableCaptures []string) {
	return b.GetAvailableMovesRaw(c, true)
}

func (b *Board) GetAvailableMovesRaw(c Coord, excludeCastling bool) (availableMoves, availableCaptures []string) {
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

func (b *Board) IsInCheck(color byte) bool {
	var targetCoord Coord
	iswhite := color == WhiteToMove
	pawnDirection := -1
	offset := uint8(0)
	king := b.GetKing(color)
	if king.Rank == 8 {
		fmt.Println(b.ExportFEN())
	}
	if iswhite {
		pawnDirection = 1
		offset = 6
	}

	targetCoord = Coord{File: king.File + 1, Rank: king.Rank + pawnDirection}
	if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == P+offset {
		return true
	}
	targetCoord = Coord{File: king.File - 1, Rank: king.Rank + pawnDirection}
	if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == P+offset {
		return true
	}

	for _, knightMove := range knightMoves {
		targetCoord = Coord{File: king.File + knightMove[0], Rank: king.Rank + knightMove[1]}
		if CoordInBounds(targetCoord) && b.AccessCoord(targetCoord) == N+offset {
			return true
		}
	}

	var n, w, s, e, ne, nw, se, sw bool = true, true, true, true, true, true, true, true
	var piece uint8
	for i := 1; i < 8; i++ {
		if n {
			n = king.Rank+i <= 7
			if n {
				piece = b.AccessCoord(Coord{king.File, king.Rank + i})
				if piece == 0 {

				} else if piece == Q+offset || piece == R+offset || (i == 1 && piece == K+offset) {
					return true
				} else {
					n = false
				}
			}
		}
		if s {
			s = king.Rank-i >= 0
			if s {
				piece = b.AccessCoord(Coord{king.File, king.Rank - i})
				if piece == 0 {

				} else if piece == Q+offset || piece == R+offset || (i == 1 && piece == K+offset) {
					return true
				} else {
					s = false
				}
			}
		}
		if w {
			w = king.File+i <= 7
			if w {
				piece = b.AccessCoord(Coord{king.File + i, king.Rank})
				if piece == 0 {

				} else if piece == Q+offset || piece == R+offset || (i == 1 && piece == K+offset) {
					return true
				} else {
					w = false
				}
			}
		}
		if e {
			e = king.File-i >= 0
			if e {
				piece = b.AccessCoord(Coord{king.File - i, king.Rank})
				if piece == 0 {

				} else if piece == Q+offset || piece == R+offset || (i == 1 && piece == K+offset) {
					return true
				} else {
					e = false
				}
			}
		}
		if nw {
			nw = king.File+i <= 7 && king.Rank+i <= 7
			if nw {
				piece = b.AccessCoord(Coord{king.File + i, king.Rank + i})
				if piece == 0 {

				} else if piece == Q+offset || piece == B+offset || (i == 1 && piece == K+offset) {
					return true
				} else {
					nw = false
				}
			}
		}
		if ne {
			ne = king.File-i >= 0 && king.Rank+i <= 7
			if ne {
				piece = b.AccessCoord(Coord{king.File - i, king.Rank + i})
				if piece == 0 {

				} else if piece == Q+offset || piece == B+offset || (i == 1 && piece == K+offset) {
					return true
				} else {
					ne = false
				}
			}
		}
		if se {
			se = king.File-i >= 0 && king.Rank-i >= 0
			if se {
				piece = b.AccessCoord(Coord{king.File - i, king.Rank - i})
				if piece == 0 {

				} else if piece == Q+offset || piece == B+offset || (i == 1 && piece == K+offset) {
					return true
				} else {
					se = false
				}
			}
		}
		if sw {
			sw = king.File+i <= 7 && king.Rank-i >= 0
			if sw {
				piece = b.AccessCoord(Coord{king.File + i, king.Rank - i})
				if piece == 0 {

				} else if piece == Q+offset || piece == B+offset || (i == 1 && piece == K+offset) {
					return true
				} else {
					sw = false
				}
			}
		}
	}
	return false
}

var pawnCaptures [2][2]int = [2][2]int{{1, 1}, {-1, 1}}

func (b *Board) GetPawnMoves(c Coord) (moves, captures []string) {
	var targetCoord Coord
	isWhite := b.Coords[c.File][c.Rank] <= PieceOffset
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

	if c.File > 0 && c.File < 7 {
		targetCoord = Coord{c.File + 1, c.Rank + direction}
		if isOpponentPiece(b, c, targetCoord) {
			captures = append(captures, CoordsToMove(c, targetCoord))
		}
		targetCoord = Coord{c.File - 1, c.Rank + direction}
		if isOpponentPiece(b, c, targetCoord) {
			captures = append(captures, CoordsToMove(c, targetCoord))
		}
	}

	if c.File == 0 {
		targetCoord = Coord{c.File + 1, c.Rank + direction}
		if isOpponentPiece(b, c, targetCoord) {
			captures = append(captures, CoordsToMove(c, targetCoord))
		}
	}

	if c.File == 7 {
		targetCoord = Coord{c.File - 1, c.Rank + direction}
		if isOpponentPiece(b, c, targetCoord) {
			captures = append(captures, CoordsToMove(c, targetCoord))
		}
	}

	if (c.Rank == 3 || c.Rank == 4) && b.EnPassantTarget != "-" {
		if c.File > 0 && c.File < 7 {
			targetCoord = Coord{c.File + 1, c.Rank + direction}
			if CoordToAlg(targetCoord) == b.EnPassantTarget {
				captures = append(captures, CoordsToMove(c, targetCoord))
			}
			targetCoord = Coord{c.File - 1, c.Rank + direction}
			if CoordInBounds(targetCoord) && CoordToAlg(targetCoord) == b.EnPassantTarget {
				captures = append(captures, CoordsToMove(c, targetCoord))
			}
		}

		if c.File == 0 {
			targetCoord = Coord{c.File + 1, c.Rank + direction}
			if CoordToAlg(targetCoord) == b.EnPassantTarget {
				captures = append(captures, CoordsToMove(c, targetCoord))
			}
		}

		if c.File == 7 {
			targetCoord = Coord{c.File - 1, c.Rank + direction}
			if CoordToAlg(targetCoord) == b.EnPassantTarget {
				captures = append(captures, CoordsToMove(c, targetCoord))
			}
		}
	}

	return
}

func (b *Board) GetBishopMoves(c Coord) (moves, captures []string) {
	var ul, ur, dl, dr bool = true, true, true, true
	var targetCoord Coord
	for i := 1; i < 8; i++ {
		if ur {
			ur = c.Rank+i <= 7 && c.File+i <= 7
			targetCoord = Coord{c.File + i, c.Rank + i}
			if ur && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if ur && isOpponentPiece(b, c, targetCoord) {
				ur = false
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				ur = false
			}
		}
		if dr {
			dr = c.Rank-i >= 0 && c.File+i <= 7
			targetCoord = Coord{c.File + i, c.Rank - i}
			if dr && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if dr && isOpponentPiece(b, c, targetCoord) {
				dr = false
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				dr = false
			}
		}

		if ul {
			ul = c.Rank+i <= 7 && c.File-i >= 0
			targetCoord = Coord{c.File - i, c.Rank + i}
			if ul && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if ul && isOpponentPiece(b, c, targetCoord) {
				ul = false
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				ul = false
			}
		}

		if dl {
			dl = c.Rank-i >= 0 && c.File-i >= 0
			targetCoord = Coord{c.File - i, c.Rank - i}
			if dl && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if dl && isOpponentPiece(b, c, targetCoord) {
				dl = false
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				dl = false
			}
		}
		if !ur && !ul && !dr && !dl {
			break
		}

	}
	return
}

var knightMoves = [8][2]int{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}}

func (b *Board) GetKnightMoves(c Coord) (moves, captures []string) {
	var targetCoord Coord

	for i := 0; i < 8; i++ {
		targetCoord = Coord{
			c.File + knightMoves[i][0], c.Rank + knightMoves[i][1],
		}
		if !CoordInBounds(targetCoord) {
			continue
		}
		if b.AccessCoord(targetCoord) == 0 {
			moves = append(moves, CoordsToMove(c, targetCoord))
		} else if isOpponentPiece(b, c, targetCoord) {
			captures = append(captures, CoordsToMove(c, targetCoord))
		}

	}

	return
}

func (b *Board) GetRookMoves(c Coord) (moves, captures []string) {
	var u, d, l, r bool = true, true, true, true
	var targetCoord Coord
	for i := 1; i < 8; i++ {
		if u {
			u = c.Rank+i <= 7
			targetCoord = Coord{c.File, c.Rank + i}
			if u && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if u && isOpponentPiece(b, c, targetCoord) {
				u = false
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				u = false
			}
		}
		if d {
			d = c.Rank-i >= 0
			targetCoord = Coord{c.File, c.Rank - i}
			if d && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if d && isOpponentPiece(b, c, targetCoord) {
				d = false
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				d = false
			}
		}

		if l {
			l = c.File-i >= 0
			targetCoord = Coord{c.File - i, c.Rank}
			if l && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if l && isOpponentPiece(b, c, targetCoord) {
				l = false
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				l = false
			}
		}

		if r {
			r = c.File+i <= 7
			targetCoord = Coord{c.File + i, c.Rank}
			if r && b.AccessCoord(targetCoord) == 0 {
				moves = append(moves, CoordsToMove(c, targetCoord))
			} else if r && isOpponentPiece(b, c, targetCoord) {
				r = false
				captures = append(captures, CoordsToMove(c, targetCoord))
			} else {
				r = false
			}
		}
		if !u && !d && !l && !r {
			break
		}

	}
	return
}

func (b *Board) GetQueenMoves(c Coord) (moves, captures []string) {
	bishopMoves, bishopCaptures := b.GetBishopMoves(c)
	rookMoves, rookCaptures := b.GetRookMoves(c)
	moves = append(bishopMoves, rookMoves...)
	captures = append(bishopCaptures, rookCaptures...)
	return
}

func (b *Board) GetKingMoves(c Coord, excludeCastling bool) (moves, captures []string) {
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

func (b *Board) isCastling(move string) bool {
	for i, castlingMove := range CastlingMoves {
		if move == castlingMove && strings.Contains(b.CastlingRights, CastlingRights[i]) {
			return true
		}
	}
	return false
}

func (b *Board) isEnPassant(move string) bool {
	from, to := move[:2], move[2:]
	piece := b.AccessCoord(AlgToCoord(from))
	return (piece == 1 || piece == 7) && to == b.EnPassantTarget
}
