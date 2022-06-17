package eval

import (
	"fmt"
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

const (
	majorDiag      float32 = 0.2 //a1d8 && a8d1
	minorDiag      float32 = 0.1 //a2g8 b1h7 && a7g1 b8h2
	knightCenter22 float32 = 0.3
	knightCenter44 float32 = 0.2
	knightOuterRim float32 = -0.2
	knightInnerRim float32 = -0.05
)

func getPieceSpecificScore(piece uint8, c board.Coord, color byte) float32 {
	switch piece {
	case board.P, board.P + 6:
		return getPawnScore(c, color)
	case board.B, board.B + 6:
		return getBishopDiagScore(c)
	case board.N, board.N + 6:
		return getKnightPositionScore(c)
	default:
		return 0
	}
}

func getPawnScore(c board.Coord, color byte) float32 {
	//add structure modifiers - doubled, passed, protected
	return getPawnAdvancementScore(c, color) + getCentralPawn(c)
}

func getPawnAdvancement(c board.Coord, color byte) int {
	switch color {
	case board.BlackToMove:
		return 6 - c.Rank
	default:
		return c.Rank - 1
	}
}

func getPawnAdvancementScore(c board.Coord, color byte) float32 {
	var score float32 = 0
	switch getPawnAdvancement(c, color) {
	case 6:
		score += 0.75
	case 5:
		score += 0.5
	case 4, 3:
		score += 0.3
	case 2:
		score += 0.2
	}

	return score
}

func getCentralPawn(c board.Coord) float32 {
	if c.Rank > 3 && c.Rank < 5 && c.File > 3 && c.File < 5 {
		return 0.2
	}
	return 0
}

func getKnightPositionScore(c board.Coord) float32 {
	if c.Rank == 0 || c.Rank == 7 || c.File == 0 || c.File == 7 {
		return knightOuterRim
	}

	if c.Rank == 1 || c.Rank == 6 || c.File == 1 || c.File == 6 {
		return knightInnerRim
	}

	if c.Rank > 3 && c.Rank < 5 && c.File > 3 && c.File < 5 {
		return knightCenter22
	}

	if c.Rank > 2 && c.Rank < 6 && c.File > 2 && c.File < 6 {
		return knightCenter44
	}

	return 0
}

func getMajorDiagScoreUR(c board.Coord) float32 {
	if c.File-c.Rank == 0 {
		return majorDiag
	}
	return 0
}

func getMajorDiagScoreDR(c board.Coord) float32 {
	if c.Rank+c.File == 7 {
		return majorDiag
	}
	return 0
}

func getMinoDiagScoreUR(c board.Coord) float32 {
	if c.Rank-c.File == 1 || c.File-c.Rank == 1 {
		return minorDiag
	}
	return 0
}

func getMinorDiagScoreDR(c board.Coord) float32 {
	if c.File+c.Rank == 6 || c.File+c.Rank == 8 {
		return minorDiag
	}

	return 0
}

func getBishopDiagScore(c board.Coord) float32 {
	return getMajorDiagScoreDR(c) + getMajorDiagScoreUR(c) + getMinoDiagScoreUR(c) + getMinorDiagScoreDR(c)
}

func SideDependantEval(e *EvalEngine, b *board.Board) float32 {
	if e.RootNode.Position.SideToMove == board.WhiteToMove {
		return GetEvaluation(e, b)
	} else {
		return GetEvaluation(e, b)
	}
}

func GetEvaluation(e *EvalEngine, b *board.Board) float32 {
	inCheck := b.IsInCheck(b.SideToMove)
	m, c := b.GetLegalMoves(b.SideToMove)
	all := append(m, c...)

	//Mate = +/-Inf score
	if inCheck && len(all) == 0 {
		if b.SideToMove == board.WhiteToMove {
			return negInf
		} else {
			return posInf
		}
		//Stale mate = 0 score
	} else if len(all) == 0 {
		return 0
	}

	whitePieces := b.GetPieces(board.WhiteToMove)
	blackPieces := b.GetPieces(board.BlackToMove)
	if DEBUG {
		fmt.Printf("white pieces %d, black pieces %d\n", len(whitePieces), len(blackPieces))
		fmt.Println("White pieces:")
	}
	var eval, pieceEval float32 = 0, 0
	for _, piece := range whitePieces {
		pieceType := b.AccessCoord(piece)
		baseWeight := pieceWeights[b.AccessCoord(piece)-1]
		moves, captures := b.GetAvailableMoves(piece)
		pieceEval = baseWeight + float32(len(moves))*moveWeight + float32(len(captures))*captureWeight + getPieceSpecificScore(pieceType, piece, board.WhiteToMove)
		if DEBUG {
			fmt.Printf("Evaluation for piece %s is %f\n", board.CoordToAlg(piece), pieceEval)
		}
		eval += pieceEval
	}

	if DEBUG {
		fmt.Println("Black pieces")
	}
	for _, piece := range blackPieces {
		pieceType := b.AccessCoord(piece)
		baseWeight := pieceWeights[b.AccessCoord(piece)-1]
		moves, captures := b.GetAvailableMoves(piece)
		pieceEval = baseWeight - float32(len(moves))*moveWeight - float32(len(captures))*captureWeight - getPieceSpecificScore(pieceType, piece, board.BlackToMove)
		if DEBUG {
			fmt.Printf("Evaluation for piece %s is %f\n", board.CoordToAlg(piece), pieceEval)
		}
		eval += pieceEval
	}

	b.IsEvaluated, b.CachedEval = true, eval
	e.Evaluations++
	return eval
}

func kingDeadScore(b *board.Board) (float32, bool) {
	whiteKingDead := !board.CoordInBounds(b.GetKing(board.WhiteToMove))
	blackKingDead := !board.CoordInBounds(b.GetKing(board.BlackToMove))
	switch {
	case whiteKingDead:
		return negInf, true
	case blackKingDead:
		return posInf, true
	default:
		return 0.0, false
	}
}
