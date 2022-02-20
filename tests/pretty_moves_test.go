package main

import (
	"github.com/likeawizard/chess-go/board"
	"testing"
)

func TestSimplePawnMove(t *testing.T) {
	var b board.Board
	move := "e2e4"
	pretty := "e4"
	b.Init()

	t.Run("Simple pawn move", func(t *testing.T) {
		if b.MoveToPretty(move) != pretty {
			t.Errorf("Got: %s Want: %s", move, pretty)
		}
	})
}

func TestPawnCapture(t *testing.T) {
	var b board.Board
	move := "e4d5"
	pretty := "exd5"
	b.ImportFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2")

	t.Run("e4 pawn takes d5 pawn", func(t *testing.T) {
		if b.MoveToPretty(move) != pretty {
			t.Errorf("Got: %s Want: %s", move, pretty)
		}
	})
}

func TestSimplePieceMoveAndRankConflict(t *testing.T) {
	var b board.Board
	move := "g3h5"
	pretty := "Nh5"
	b.ImportFEN("rnbqkbnr/pppppppp/8/8/8/2N3N1/PPPPPPPP/R1BQKB1R w KQkq - 0 1")

	t.Run("Knight moves to h5 no conflict", func(t *testing.T) {
		prettified := b.MoveToPretty(move)
		if prettified != pretty {
			t.Errorf("Got: %s Want: %s", prettified, pretty)
		}
	})

	move = "c3e4"
	pretty = "Nce4"

	t.Run("Knight from c3 moves to no conflict", func(t *testing.T) {
		prettified := b.MoveToPretty(move)
		if prettified != pretty {
			t.Errorf("Got: %s Want: %s", prettified, pretty)
		}
	})
}
