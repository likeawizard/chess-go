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
	initPieceWeightLUT()
}

type pieceEvalFn func(*board.Board, board.Square, int) int

var pieceEvals = [6]pieceEvalFn{pawnEval, bishopEval, knightEval, rookEval, queenEval, kingEval}

func pawnEval(b *board.Board, sq board.Square, side int) int {
	value := 0
	if IsProtected(b, sq, side) {
		value += weights.Pawn.Protected
	}
	if IsDoubled(b, sq, side) {
		value += weights.Pawn.Doubled
	}

	if IsIsolated(b, sq, side) {
		value += weights.Pawn.Isolated
	}
	advancmentValue := weights.Pawn.Advance
	if IsPassed(b, sq, side) {
		advancmentValue += weights.Pawn.Passed
	}

	value += getCentralPawn(sq)

	value += advancmentValue * getPawnAdvancement(sq, side)

	return value
}

func queenEval(b *board.Board, sq board.Square, side int) int {
	return 0
}

// TODO: combine all pawn functions in one with multi value return
// Piece protected a pawn
func IsProtected(b *board.Board, sq board.Square, side int) bool {
	return board.PawnAttacks[side^1][sq]^b.Pieces[side][board.PAWNS] != 0
}

// TODO: create a lookup table for files to avoid branching
func IsDoubled(b *board.Board, sq board.Square, side int) bool {
	file := sq % 8
	fileMask := board.BBoard(0)
	switch file {
	case 0:
		fileMask = board.AFile
	case 1:
		fileMask = board.BFile
	case 2:
		fileMask = board.CFile
	case 3:
		fileMask = board.DFile
	case 4:
		fileMask = board.EFile
	case 5:
		fileMask = board.FFile
	case 6:
		fileMask = board.GFile
	case 7:
		fileMask = board.HFile
	}

	return (b.Pieces[side][board.PAWNS] & fileMask).Count() > 1
}

// Has no friendly pawns on neighboring files
func IsIsolated(b *board.Board, sq board.Square, side int) bool {
	file := sq % 8
	fileMask := board.BBoard(0)
	switch file {
	case 0:
		fileMask = board.BFile
	case 1:
		fileMask = board.AFile | board.CFile
	case 2:
		fileMask = board.DFile | board.BFile
	case 3:
		fileMask = board.EFile | board.CFile
	case 4:
		fileMask = board.DFile | board.FFile
	case 5:
		fileMask = board.EFile | board.GFile
	case 6:
		fileMask = board.FFile | board.HFile
	case 7:
		fileMask = board.GFile
	}

	return (b.Pieces[side][board.PAWNS] & fileMask).Count() > 0
}

// Has no opponent opposing pawns in front (same or neighbor files)
// TODO: stub
func IsPassed(b *board.Board, sq board.Square, side int) bool {
	return false
}

func getPawnAdvancement(c board.Square, side int) int {
	if side == board.WHITE {
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

func knightEval(b *board.Board, sq board.Square, side int) int {
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

func bishopPairEval(b *board.Board, side int) int {
	if b.Pieces[side][board.BISHOPS].Count() == 2 {
		return 50
	}
	return 0
}

func bishopEval(b *board.Board, sq board.Square, side int) int {
	return bishopPairEval(b, side) + getMajorDiagScoreDR(sq) + getMajorDiagScoreUR(sq) + getMinoDiagScoreUR(sq) + getMinorDiagScoreDR(sq)
}

func (e *EvalEngine) GetEvaluation(b *board.Board) int {
	e.Stats.evals++
	inCheck := b.IsChecked(b.Side)
	all := b.MoveGen()

	//Mate = +/-Inf score
	if inCheck && len(all) == 0 {
		if b.Side == board.WHITE {
			return -math.MaxInt
		} else {
			return math.MaxInt
		}
		//Stalemate = 0 score
	} else if len(all) == 0 {
		return 0
	}

	var eval, pieceEval int = 0, 0
	phase := getGamePhase(b)

	// TODO: ensure no move gen is dependent on b.IsWhite internally
	var side = -1
	for color := board.WHITE; color <= board.BLACK; color++ {
		side *= -1
		numPawns := b.Pieces[color][board.PAWNS].Count()
		for pieceType := board.PAWNS; pieceType <= board.KINGS; pieceType++ {
			pieces := b.Pieces[color][pieceType]
			for pieces > 0 {
				piece := pieces.PopLS1B()
				pieceEval = PieceWeights[pieceType] + PiecePawnBonus[pieceType][numPawns]
				// Tapered eval - more bias towards PST in the opening and more bias to individual eval functions towards the endgame
				pieceEval += (PST[color][pieceType][piece]*(256-phase) + pieceEvals[pieceType](b, board.Square(piece), color)*phase) / 256
				// moves := b.GetMovesForPiece(board.Square(piece), pieceType, 0, 0)
				// pieceEval += + len(moves)*weights.Moves.Move
				eval += side * pieceEval
			}
		}
	}

	return eval
}

// Determine the game phase as a sliding factor between opening and endgame
// https://www.chessprogramming.org/Tapered_Eval#Implementation_example
func getGamePhase(b *board.Board) (phase int) {
	phase = 24

	for color := board.WHITE; color <= board.BLACK; color++ {
		for pieceType := board.PAWNS; pieceType <= board.KINGS; pieceType++ {
			switch pieceType {
			case board.BISHOPS, board.KINGS:
				phase -= b.Pieces[color][pieceType].Count()
			case board.ROOKS:
				phase -= 2 * b.Pieces[color][pieceType].Count()
			case board.QUEENS:
				phase -= 4 * b.Pieces[color][pieceType].Count()
			}
		}
	}

	phase = (phase * 268) / 24

	return
}

// TODO: try branchless: eliminate min/max and use branchless abs()
func distCenter(sq board.Square) int {
	c := int(sq)
	return Max(3-c/8, c/8-4) + Max(3-c%8, c%8-4)
}

func distSqares(us, them board.Square) int {
	u, t := int(us), int(them)
	return Max((u-t)/8, (t-u)/8) + Max((u-t)%8, (t-u)%8)
}

// King safety score as a measure of distance from the board center and the number of adjacent enemy pieces and friendly pieces
// TODO: naive initial approach
// use board direction to value own pieces: i.e. a King in front of 3 pawns is not the same as a king behind 3 pawns
// consider using opponent piece attacks around the king instead of actual pieces. use piece weights for opponent threat levels: a queen near our king should be a larger concern than a bishop
func getKingSafety(b *board.Board, king board.Square, side int) (kingSafety int) {
	kingSafety += 2 * distCenter(king)

	numFriendly := board.KingAttacks[king] & b.Occupancy[side]
	numOpponent := board.KingAttacks[king] & b.Occupancy[side^1]

	kingSafety += 5*numFriendly.Count() - 15*numOpponent.Count()

	return
}

func getKingActivity(b *board.Board, king board.Square, side int) (kingActivity int) {
	oppKing := b.Pieces[b.Side^1][board.KINGS].LS1B()
	kingActivity = -(distCenter(king) + distSqares(king, board.Square(oppKing)))
	return
}

func kingEval(b *board.Board, king board.Square, side int) int {
	phase := getGamePhase(b)
	return (getKingSafety(b, king, side)*(256-phase) + getKingActivity(b, king, side)*phase) / 256

}

// Evaluation for rooks - connected & (semi)open files
// TODO: stub
func rookEval(b *board.Board, rook board.Square, side int) (rookScore int) {
	return
	// offset := uint8(6)
	// if side == board.WHITE {
	// 	offset = 0
	// }
	// hasOwnPawn, hasOppPawn, connected := false, false, false
	// for dirIdx := 0; dirIdx < 4; dirIdx++ {
	// 	for i := board.Square(1); i <= board.CompassBlock[rook][dirIdx]; i++ {
	// 		target := rook + i*board.Compass[dirIdx]

	// 		if b.Coords[target] == 0 {
	// 			continue
	// 		}
	// 		if b.Coords[target] == board.R+offset {
	// 			connected = true
	// 		}

	// 		//Only look at N and S
	// 		if dirIdx > 2 {
	// 			continue
	// 		}

	// 		//Check for own or opponent pawns
	// 		switch b.Coords[target] {
	// 		case 7 - offset:
	// 			hasOppPawn = true
	// 		case 1 + offset:
	// 			hasOwnPawn = true
	// 		}
	// 	}
	// }

	// if !hasOwnPawn && !hasOppPawn {
	// 	rookScore += 15
	// } else if hasOppPawn && !hasOwnPawn {
	// 	rookScore += 10
	// }

	// if connected {
	// 	rookScore += 10
	// }

	// return
}
