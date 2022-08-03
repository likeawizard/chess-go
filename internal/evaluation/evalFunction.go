package eval

import (
	"fmt"
	"math"

	"github.com/likeawizard/chess-go/internal/board"
)

var weights *Weights

func init() {
	var err error
	weights, err = LoadWeights()
	if err != nil {
		fmt.Println("Unable to load weights")
		panic(1)
	}
	fmt.Printf("%+v\n", weights)
}

var (
	PieceWeights = [12]float32{1, 3.2, 2.9, 5, 9, 0, -1, -3.2, -2.9, -5, -9, 0}
)

func getPieceSpecificScore(b *board.Board, piece uint8, c board.Square, isWhite bool) int {
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

func getPawnScore(b *board.Board, sq board.Square, isWhite bool) (value int) {
	value = 0
	if IsProtected(b, sq, isWhite) {
		value += weights.Pawn.Protected
	}
	if IsDoubled(b, sq, isWhite) {
		value += weights.Pawn.Doubled
	}

	if IsIsolated(b, sq, isWhite) {
		value += weights.Pawn.Isolated
	}
	advancmentValue := weights.Pawn.Advance
	if IsPassed(b, sq, isWhite) {
		advancmentValue += weights.Pawn.Passed
	}

	value += getCentralPawn(sq)

	value += advancmentValue * getPawnAdvancement(sq, isWhite)

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

func getPawnAdvancement(c board.Square, isWhite bool) int {
	if isWhite {
		return int(c/8 - 1)
	} else {
		return int(6 - c/8)
	}
}

func getCentralPawn(sq board.Square) int {
	switch {
	case (sq/8 == 3 || sq/8 == 4) && (sq%8 == 3 || sq%8 == 4):
		return weights.Pawn.Center22
	case (sq/8 == 2 || sq/8 == 5) && (sq%8 == 2 || sq%8 == 5):
		return weights.Pawn.Center44
	default:
		return 0
	}
}

func getKnightPositionScore(sq board.Square) int {
	switch {
	case (sq/8 == 3 || sq/8 == 4) && (sq%8 == 3 || sq%8 == 4):
		return weights.Knight.Center22
	case (sq/8 == 2 || sq/8 == 5) && (sq%8 == 2 || sq%8 == 5):
		return weights.Knight.Center44
	case (sq/8 == 1 || sq/8 == 6) && (sq%8 == 1 || sq%8 == 6):
		return weights.Knight.InnerRim
	default:
		return weights.Knight.OuterRim
	}

}

func getMajorDiagScoreUR(c board.Square) int {
	if c%9 == 0 {
		return weights.Bishop.MajorDiag
	}
	return 0
}

func getMajorDiagScoreDR(c board.Square) int {
	if c%7 == 0 {
		return weights.Bishop.MajorDiag
	}
	return 0
}

func getMinoDiagScoreUR(c board.Square) int {
	if c%9 == 1 || c%9 == 8 {
		return weights.Bishop.MinorDiag
	}
	return 0
}

func getMinorDiagScoreDR(c board.Square) int {
	if c%7 == 6 || c%7 == 1 {
		return weights.Bishop.MinorDiag
	}

	return 0
}

func getBishopDiagScore(c board.Square) int {
	return getMajorDiagScoreDR(c) + getMajorDiagScoreUR(c) + getMinoDiagScoreUR(c) + getMinorDiagScoreDR(c)
}

func GetEvaluation(e *EvalEngine, b *board.Board) int {
	inCheck := b.IsInCheck(b.IsWhite)
	all := b.GetLegalMoves()

	//Mate = +/-Inf score
	if inCheck && len(all) == 0 {
		if b.IsWhite {
			return -math.MaxInt
		} else {
			return math.MaxInt
		}
		//Stale mate = 0 score
	} else if len(all) == 0 {
		return 0
	}

	whitePieces := b.GetPieces(true)
	blackPieces := b.GetPieces(false)
	var eval, pieceEval int = 0, 0

	// TODO: ensure no move gen is dependent on b.IsWhite internally
	isWhite := b.IsWhite
	b.IsWhite = true
	for _, piece := range whitePieces {
		pieceVal := b.Coords[piece]
		// TODO: eval for pinned pieces?
		moves := b.GetMovesForPiece(piece, 0, 0)
		pieceEval = getPieceWeight(pieceVal) + len(moves)*weights.Moves.Move + getPieceSpecificScore(b, pieceVal, piece, true)
		eval += pieceEval
	}

	b.IsWhite = false
	for _, piece := range blackPieces {
		pieceVal := b.Coords[piece]
		moves := b.GetMovesForPiece(piece, 0, 0)
		pieceEval = getPieceWeight(pieceVal) + len(moves)*weights.Moves.Move + getPieceSpecificScore(b, pieceVal, piece, false)
		eval -= pieceEval
	}

	b.IsWhite = isWhite
	e.Evaluations++
	return eval
}
