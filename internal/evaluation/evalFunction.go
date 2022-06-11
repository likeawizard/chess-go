package eval

import (
	"math"

	"github.com/likeawizard/chess-go/internal/board"
)

func SideDependantEval(e *EvalEngine, b *board.Board) float64 {
	if e.RootNode.Position.SideToMove == board.WhiteToMove {
		return GetEvaluation2(e, b)
	} else {
		return GetEvaluation2(e, b)
	}
}

func GetEvaluation(e *EvalEngine, b *board.Board) float64 {
	const (
		moveWeight    = 0.02
		captureWeight = 0.04
	)
	var pieceWeights = [12]float64{1, 3.2, 2.9, 5, 9, 0, -1, -3.2, -2.9, -5, -9, 0}
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
	allPieces := append(whitePieces, blackPieces...)
	var eval = 0.0
	for _, piece := range allPieces {
		baseWeight := pieceWeights[b.AccessCoord(piece)-1]
		moves, captures := b.GetAvailableMoves(piece)
		eval += baseWeight * (1.0 + float64(len(moves))*moveWeight + float64(len(captures))*captureWeight)
	}

	b.IsEvaluated, b.CachedEval = true, eval
	e.Evaluations++
	return eval
}

func GetEvaluation2(e *EvalEngine, b *board.Board) float64 {
	const (
		moveWeight    = 0.02
		captureWeight = 0.04
	)
	var pieceWeights = [12]float64{1, 3.2, 2.9, 5, 9, 0, -1, -3.2, -2.9, -5, -9, 0}
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
	var eval = 0.0
	for _, piece := range whitePieces {
		baseWeight := pieceWeights[b.AccessCoord(piece)-1]
		moves, captures := b.GetAvailableMoves(piece)
		eval += baseWeight + getSign(board.WhiteToMove)*(float64(len(moves))*moveWeight+float64(len(captures))*captureWeight)
	}

	for _, piece := range blackPieces {
		baseWeight := pieceWeights[b.AccessCoord(piece)-1]
		moves, captures := b.GetAvailableMoves(piece)
		eval += baseWeight + getSign(board.BlackToMove)*(float64(len(moves))*moveWeight+float64(len(captures))*captureWeight)
	}

	b.IsEvaluated, b.CachedEval = true, eval
	e.Evaluations++
	return eval
}

func getSign(side string) float64 {
	switch side {
	case board.WhiteToMove:
		return 1
	case board.BlackToMove:
		return -1
	}
	return 1
}

func kingDeadScore(b *board.Board) (float64, bool) {
	whiteKingDead := !board.CoordInBounds(board.GetKing(*b, board.WhiteToMove))
	blackKingDead := !board.CoordInBounds(board.GetKing(*b, board.BlackToMove))
	switch {
	case whiteKingDead:
		return math.Inf(-1), true
	case blackKingDead:
		return math.Inf(1), true
	default:
		return 0.0, false
	}
}
