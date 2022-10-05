package eval

import (
	"context"
	"fmt"
	"sort"

	"github.com/likeawizard/chess-go/internal/board"
	"github.com/likeawizard/chess-go/internal/book"
	"github.com/likeawizard/chess-go/internal/config"
)

type PickBookMove func(*board.Board) board.Move

type EvalEngine struct {
	Stats        EvalStats
	Board        *board.Board
	EnableBook   bool
	PickBookMove PickBookMove
	SearchDepth  int
	TTable       *TTable
}

func NewEvalEngine(b *board.Board, c *config.Config) (*EvalEngine, error) {
	var bookMethod PickBookMove = book.GetWeighted
	switch c.Book.Method {
	case "best":
		bookMethod = book.GetBest
	case "weighted":
		bookMethod = book.GetWeighted
	}
	return &EvalEngine{
		Board:        b,
		EnableBook:   c.Book.Enable,
		PickBookMove: bookMethod,
		SearchDepth:  c.Engine.MaxDepth,
		TTable:       NewTTable(c.Engine.TTSize),
	}, nil
}

// Returns the best move and best opponent response - ponder
func (e *EvalEngine) GetMove(ctx context.Context, pv *[]board.Move, silent bool) (board.Move, board.Move) {
	var best, ponder board.Move
	var ok bool
	all := e.Board.MoveGen()
	if len(all) == 1 {
		best = all[0]
	} else if e.EnableBook && book.InBook(e.Board) {
		book.PrintBookMoves(e.Board)
		move := e.PickBookMove(e.Board)
		fmt.Println("Picking Book move: ", move)
		return move, 0
	} else {
		best, ponder, ok = e.IDSearch(ctx, e.SearchDepth, pv, silent)
		if !ok {
			best = all[0]
		}
	}

	return best, ponder
}

func (e *EvalEngine) OrderMoves(pv board.Move, moves *[]board.Move) {
	sort.Slice(*moves, func(i int, j int) bool {
		return (*moves)[i] == pv || e.getMoveValue((*moves)[i]) > e.getMoveValue((*moves)[j])
	})
}

// Estimate the potential strength of the move for move ordering
func (e *EvalEngine) getMoveValue(move board.Move) (value int) {

	if move.IsCapture() {
		attacker := PieceWeights[(move.Piece()-1)%6]
		_, _, victim := e.Board.PieceAtSquare(move.To())
		value = PieceWeights[victim] - attacker/2
	}

	// TODO: implement SEE or MVV-LVA ordering
	// Calculate the relative value of exchange
	// from, to := move.FromTo()
	// us, them := PieceWeights[b.Coords[from]], PieceWeights[b.Coords[to]]
	// if them == 0 {
	// 	value += 0
	// } else {
	// 	value += dir * (0.5*us + them)
	// }

	// Prioritize promotions
	if move.Promotion() != 0 {
		value += 3
	}

	return
}

// TODO: try branchless optimization
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
