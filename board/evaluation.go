package board

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

var DEBUG = false

const (
	moveWeight    = 0.02
	captureWeight = 0.04
)

var pieceWeights = [12]float64{1, 3.2, 2.9, 5, 9, 0, -1, -3.2, -2.9, -5, -9, 0}
var Evaluations int64
var CachedEvals int64

var rootEvalNode *Line
var wg sync.WaitGroup

type Line struct {
	Position   *Board
	Parent     *Line
	Candidates []*Line
}

func InitEvalEngine(b *Board) {
	rootEvalNode = &Line{
		Position: b,
	}
	DEBUG, _ = strconv.ParseBool(os.Getenv("EVALUATION_DEBUG"))
}

func (l *Line) GetEval() float64 {
	return l.Position.GetEvaluation()
}

func (l *Line) buildChildren() (moves []string) {
	fen := l.Position.ExportFEN()
	m, c := l.Position.GetMoves(l.Position.sideToMove)
	all := append(c, m...)
	candidateMoves := make([]*Line, len(all))
	for i := 0; i < len(candidateMoves); i++ {
		candidateMoves[i] = &Line{}
		candidateMoves[i].Parent = l
		candidateMoves[i].Position = &Board{}
		candidateMoves[i].Position.ImportFEN(fen)
		candidateMoves[i].Position.MoveLongAlg(all[i])
	}

	l.Candidates = candidateMoves

	return all
}

func (l *Line) minmax(depth int) float64 {
	if depth == 0 {
		return l.Position.GetEvaluation()
	}

	if l.Candidates == nil {
		l.buildChildren()
	}

	if l.Position.sideToMove == whiteToMove {
		val := math.Inf(-1)
		for i := 0; i < len(l.Candidates); i++ {
			val = math.Max(val, (l.Candidates)[i].minmax(depth-1))
		}
		return val
	} else {
		val := math.Inf(1)
		for i := 0; i < len(l.Candidates); i++ {
			val = math.Min(val, l.Candidates[i].minmax(depth-1))
		}
		return val
	}
}

func (l *Line) alphabeta(depth int, alpha, beta float64) float64 {
	if depth == 0 {
		return l.Position.GetEvaluation()
	}

	if l.Candidates == nil {
		l.buildChildren()
	}

	if l.Position.sideToMove == whiteToMove {
		val := math.Inf(-1)
		for i := 0; i < len(l.Candidates); i++ {
			val = math.Max(val, l.Candidates[i].alphabeta(depth-1, alpha, beta))
			alpha = math.Max(alpha, val)
			if val >= beta {
				break
			}
		}
		return val
	} else {
		val := math.Inf(1)
		for i := 0; i < len(l.Candidates); i++ {
			val = math.Min(val, l.Candidates[i].alphabeta(depth-1, alpha, beta))
			beta = math.Min(beta, val)
			if val <= alpha {
				break
			}
		}
		return val
	}
}

func (b *Board) GetMove(depth int) string {
	if rootEvalNode.Candidates == nil {
		rootEvalNode.buildChildren()
	}
	captures, moves := b.GetMoves(b.sideToMove)
	candidateMoves := append(captures, moves...)

	moveStrengths := make([]float64, len(rootEvalNode.Candidates))
	bestMoveIndex := 0
	wg.Add(len(moveStrengths))
	for i := 0; i < len(moveStrengths); i++ {
		ii := i
		go func() {
			defer wg.Done()
			//moveStrengths[ii] = rootEvalNode.Candidates[ii].minmax(depth)
			moveStrengths[ii] = rootEvalNode.Candidates[ii].alphabeta(depth, math.Inf(-1), math.Inf(1))
			//moveStrengths[ii] = line.Candidates[ii].alphabeta(depth, -15, 15)
		}()
	}

	wg.Wait()
	if b.sideToMove == whiteToMove {
		bestMove := math.Inf(-1)
		for i := 0; i < len(moveStrengths); i++ {
			if moveStrengths[i] > bestMove {
				bestMove = moveStrengths[i]
				bestMoveIndex = i
			}
		}
	} else {
		bestMove := math.Inf(1)
		for i := 0; i < len(moveStrengths); i++ {
			if moveStrengths[i] < bestMove {
				bestMove = moveStrengths[i]
				bestMoveIndex = i
			}
		}
	}
	rootEvalNode = rootEvalNode.Candidates[bestMoveIndex]
	if DEBUG {
		for i := 0; i < len(candidateMoves); i++ {
			fmt.Printf("(%s %.2f) ", candidateMoves[i], moveStrengths[i])
		}
		fmt.Println()
	}
	return candidateMoves[bestMoveIndex]
}

func (b *Board) GetEvaluation() float64 {
	if b.isEvaluated {
		CachedEvals++
		return b.cachedEval
	}

	whiteKingDead := !CoordInBounds(GetKing(*b, whiteToMove))
	blackKingDead := !CoordInBounds(GetKing(*b, blackToMove))
	switch {
	case whiteKingDead:
		return math.Inf(-1)
	case blackKingDead:
		return math.Inf(1)
	}
	whitePieces := GetPieces(*b, whiteToMove)
	blackPieces := GetPieces(*b, blackToMove)
	allPieces := append(whitePieces, blackPieces...)
	var eval = 0.0
	for _, piece := range allPieces {
		baseWeight := pieceWeights[b.AccessCoord(piece)-1]
		moves, captures := b.GetAvailableMoves(piece)
		eval += baseWeight * (1.0 + float64(len(moves))*moveWeight + float64(len(captures))*captureWeight)
		eval += baseWeight
	}

	b.isEvaluated, b.cachedEval = true, eval
	Evaluations++
	return eval
}
