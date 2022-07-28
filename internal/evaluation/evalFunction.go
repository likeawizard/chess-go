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

func getPieceSpecificScore(piece uint8, c board.Square, color byte) float32 {
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

func getPawnScore(c board.Square, color byte) float32 {
	//add structure modifiers - doubled, passed, protected
	return getPawnAdvancementScore(c, color) + getCentralPawn(c)
}

func getPawnAdvancement(c board.Square, color byte) board.Square {
	switch color {
	case board.BlackToMove:
		return 6 - c/8
	default:
		return c/8 - 1
	}
}

func getPawnAdvancementScore(c board.Square, color byte) float32 {
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

func getCentralPawn(c board.Square) float32 {
	switch c {
	case 27, 28, 35, 36:
		return 0.2
	default:
		return 0
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
	m, c := b.GetLegalMoves(b.IsWhite)
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
	if DEBUG {
		fmt.Printf("white pieces %d, black pieces %d\n", len(whitePieces), len(blackPieces))
		fmt.Println("White pieces:")
	}
	var eval, pieceEval float32 = 0, 0
	for _, piece := range whitePieces {
		pieceType := b.Coords[piece]
		baseWeight := pieceWeights[b.Coords[piece]-1]
		moves, captures := b.GetAvailableMoves(piece)
		pieceEval = baseWeight + float32(len(moves))*moveWeight + float32(len(captures))*captureWeight + getPieceSpecificScore(pieceType, piece, board.WhiteToMove)
		if DEBUG {
			fmt.Printf("Evaluation for piece %s is %f\n", piece, pieceEval)
		}
		eval += pieceEval
	}

	if DEBUG {
		fmt.Println("Black pieces")
	}
	for _, piece := range blackPieces {
		pieceType := b.Coords[piece]
		baseWeight := pieceWeights[b.Coords[piece]-1]
		moves, captures := b.GetAvailableMoves(piece)
		pieceEval = baseWeight - float32(len(moves))*moveWeight - float32(len(captures))*captureWeight - getPieceSpecificScore(pieceType, piece, board.BlackToMove)
		if DEBUG {
			fmt.Printf("Evaluation for piece %s is %f\n", piece, pieceEval)
		}
		eval += pieceEval
	}

	e.Evaluations++
	return eval
}
