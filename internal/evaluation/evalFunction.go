package eval

import (
	"math"

	"github.com/likeawizard/chess-go/internal/board"
)

const (
	moveWeight    float32 = 0.02
	captureWeight float32 = 0.04
)

var (
	negInf       float32 = -math.MaxFloat32
	posInf       float32 = math.MaxFloat32
	PieceWeights         = [12]float32{1, 3.2, 2.9, 5, 9, 0, -1, -3.2, -2.9, -5, -9, 0}
)

const (
	majorDiag      float32 = 0.2 //a1d8 && a8d1
	minorDiag      float32 = 0.1 //a2g8 b1h7 && a7g1 b8h2
	knightCenter22 float32 = 0.3
	knightCenter44 float32 = 0.2
	knightOuterRim float32 = -0.2
	knightInnerRim float32 = -0.05
)

func getPieceSpecificScore(b *board.Board, piece uint8, c board.Square, isWhite bool) float32 {
	switch piece {
	case board.P, board.P + 6:
		return getPawnScore(b, c, isWhite)
	case board.B, board.B + 6:
		return getBishopDiagScore(c)
	case board.N, board.N + 6:
		return getKnightPositionScore(c)
	default:
		return 0
	}
}

const (
	isolated      float32 = -0.2
	doubled       float32 = -0.2
	protected     float32 = 0.3
	passedPerRank float32 = 0.1
)

func getPawnScore(b *board.Board, sq board.Square, isWhite bool) (value float32) {
	value = 1
	if IsProtected(b, sq, isWhite) {
		value += 0.1
	}
	if IsDoubled(b, sq, isWhite) {
		value -= 0.2
	}

	if IsIsolated(b, sq, isWhite) {
		value -= 0.1
	}
	advancmentValue := float32(0.075)
	if IsPassed(b, sq, isWhite) {
		advancmentValue = 0.11
	}

	if getCentralPawn(sq) {
		value += 0.3
	}

	value += advancmentValue * float32(getPawnAdvancement(sq, isWhite))

	return
}

// TODO: combine all pawn functions in one with multi value return
// Piece protected a pawn
func IsProtected(b *board.Board, sq board.Square, isWhite bool) bool {
	direction := board.Square(8)
	pawn := uint8(7)

	if isWhite {
		direction = -8
		pawn = 1
	}

	target := sq + direction + 1
	if sq%8 != 7 && b.Coords[target] == pawn {
		return true
	}

	target = sq + direction - 1
	if sq%8 != 0 && b.Coords[target] == pawn {
		return true
	}

	return false
}

func IsDoubled(b *board.Board, sq board.Square, isWhite bool) bool {
	pawn := uint8(7)
	if isWhite {
		pawn = 1
	}
	file := sq % 8
	for rank := board.Square(0); rank < 8; rank++ {
		target := file + rank*8
		if target != sq && b.Coords[target] == pawn {
			return true
		}
	}
	return false
}

// Has no friendly pawns on neighboring files
func IsIsolated(b *board.Board, sq board.Square, isWhite bool) bool {
	pawn := uint8(7)
	if isWhite {
		pawn = 1
	}
	file := sq % 8
	for rank := board.Square(0); rank < 8; rank++ {
		target := file + rank*8
		if file != 7 && b.Coords[target+1] == pawn {
			return false
		}

		if file != 0 && b.Coords[target-1] == pawn {
			return false
		}
	}
	return true
}

// Has no opponent opposing pawns in front (same or neighbor files)
func IsPassed(b *board.Board, sq board.Square, isWhite bool) bool {
	pawn := uint8(1)
	direction := board.Square(-1)
	if isWhite {
		direction = 1
		pawn = 7
	}
	file := sq % 8
	for rank := (sq / 8) + direction; rank < 8 && rank > 0; rank += direction {
		target := file + rank*8
		if b.Coords[target] == pawn {
			return false
		}

		if file != 7 && b.Coords[target+1] == pawn {
			return false
		}

		if file != 0 && b.Coords[target-1] == pawn {
			return false
		}
	}
	return true
}

func getPawnAdvancement(c board.Square, isWhite bool) board.Square {
	if isWhite {
		return c/8 - 1
	} else {
		return 6 - c/8
	}
}

func getCentralPawn(c board.Square) bool {
	switch c {
	case 27, 28, 35, 36:
		return true
	default:
		return false
	}
}

func getKnightPositionScore(c board.Square) float32 {
	switch c {
	case 27, 28, 35, 36:
		return knightCenter22
	case 42, 43, 44, 45, 34, 37, 26, 29, 18, 19, 20, 21:
		return knightCenter44
	case 49, 50, 51, 52, 53, 54, 9, 10, 11, 12, 13, 14, 41, 33, 25, 17, 46, 38, 30, 22:
		return knightInnerRim
	default:
		return knightOuterRim
	}

}

func getMajorDiagScoreUR(c board.Square) float32 {
	if c%9 == 0 {
		return majorDiag
	}
	return 0
}

func getMajorDiagScoreDR(c board.Square) float32 {
	if c%7 == 0 {
		return majorDiag
	}
	return 0
}

func getMinoDiagScoreUR(c board.Square) float32 {
	if c%9 == 1 || c%9 == 8 {
		return minorDiag
	}
	return 0
}

func getMinorDiagScoreDR(c board.Square) float32 {
	if c%7 == 6 || c%7 == 1 {
		return minorDiag
	}

	return 0
}

func getBishopDiagScore(c board.Square) float32 {
	return getMajorDiagScoreDR(c) + getMajorDiagScoreUR(c) + getMinoDiagScoreUR(c) + getMinorDiagScoreDR(c)
}

func GetEvaluation(e *EvalEngine, b *board.Board) float32 {
	inCheck := b.IsInCheck(b.IsWhite)
	m, c := b.GetLegalMoves()
	all := append(m, c...)

	//Mate = +/-Inf score
	if inCheck && len(all) == 0 {
		if b.IsWhite {
			return negInf
		} else {
			return posInf
		}
		//Stale mate = 0 score
	} else if len(all) == 0 {
		return 0
	}

	whitePieces := b.GetPieces(true)
	blackPieces := b.GetPieces(false)
	var eval, pieceEval float32 = 0, 0
	for _, piece := range whitePieces {
		pieceType := b.Coords[piece]
		baseWeight := PieceWeights[b.Coords[piece]-1]
		// TODO: eval for pinned pieces?
		moves, captures := b.GetMovesForPiece(piece, 0, 0)
		pieceEval = baseWeight + float32(len(moves))*moveWeight + float32(len(captures))*captureWeight + getPieceSpecificScore(b, pieceType, piece, true)
		eval += pieceEval
	}

	for _, piece := range blackPieces {
		pieceType := b.Coords[piece]
		baseWeight := PieceWeights[b.Coords[piece]-1]
		moves, captures := b.GetMovesForPiece(piece, 0, 0)
		pieceEval = baseWeight - float32(len(moves))*moveWeight - float32(len(captures))*captureWeight - getPieceSpecificScore(b, pieceType, piece, false)
		eval += pieceEval
	}

	e.Evaluations++
	return eval
}
