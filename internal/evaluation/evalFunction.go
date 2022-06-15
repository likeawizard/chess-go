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
	pieceWeights         = [12]float32{1, 3.2, 2.9, 5, 9, 0, -1, -3.2, -2.9, -5, -9, 0}
)

func SideDependantEval(e *EvalEngine, b *board.Board) float32 {
	if e.RootNode.Position.SideToMove == board.WhiteToMove {
		return GetEvaluation(e, b)
	} else {
		return GetEvaluation(e, b)
	}
}

func GetEvaluation(e *EvalEngine, b *board.Board) float32 {
	if b.IsEvaluated {
		CachedEvals++
		return b.CachedEval
	}

	kingDeadScore, isDead := kingDeadScore(b)

	if isDead {
		return kingDeadScore
	}

	whitePieces := board.GetPieces(*b, board.WhiteToMove)
	blackPieces := board.GetPieces(*b, board.BlackToMove)
	var eval float32 = 0.0
	for _, piece := range whitePieces {
		baseWeight := pieceWeights[b.AccessCoord(piece)-1]
		moves, captures := b.GetAvailableMoves(piece)
		eval += baseWeight + (float32(len(moves))*moveWeight + float32(len(captures))*captureWeight)
	}

	for _, piece := range blackPieces {
		baseWeight := pieceWeights[b.AccessCoord(piece)-1]
		moves, captures := b.GetAvailableMoves(piece)
		eval += baseWeight - (float32(len(moves))*moveWeight + float32(len(captures))*captureWeight)
	}

	b.IsEvaluated, b.CachedEval = true, eval
	e.Evaluations++
	return eval
}

func kingDeadScore(b *board.Board) (float32, bool) {
	whiteKingDead := !board.CoordInBounds(board.GetKing(*b, board.WhiteToMove))
	blackKingDead := !board.CoordInBounds(board.GetKing(*b, board.BlackToMove))
	switch {
	case whiteKingDead:
		return negInf, true
	case blackKingDead:
		return posInf, true
	default:
		return 0.0, false
	}
}
