package board

import (
	"testing"

	"github.com/likeawizard/chess-go/internal/board"
)

func TestBishopCheck(t *testing.T) {
	var b board.Board
	expected := true
	b.ImportFEN("bk6/8/8/8/8/8/8/7K w - - 0 1")

	t.Run("Bishop Check", func(t *testing.T) {
		if b.IsInCheck(board.WhiteToMove) != expected {
			t.Errorf("Check not detected")
		}
	})
	b.ImportFEN("bk6/8/8/8/8/8/8/7K w - - 0 1")

	t.Run("Bishop Check", func(t *testing.T) {
		if b.IsInCheck(board.WhiteToMove) != expected {
			t.Errorf("Check not detected")
		}
	})

	b.ImportFEN("1k5K/8/8/8/8/2b5/8/8 w - - 0 1")

	t.Run("Bishop Check", func(t *testing.T) {
		if b.IsInCheck(board.WhiteToMove) != expected {
			t.Errorf("Check not detected")
		}
	})

	b.ImportFEN("1k5K/8/8/8/8/8/7B/8 b - - 0 1")

	t.Run("Bishop Check", func(t *testing.T) {
		if b.IsInCheck(board.BlackToMove) != expected {
			t.Errorf("Check not detected")
		}
	})

	b.ImportFEN("2k4K/8/8/8/8/8/7B/8 b - - 0 1")

	t.Run("Bishop Check", func(t *testing.T) {
		if b.IsInCheck(board.BlackToMove) == expected {
			t.Errorf("False check detected")
		}
	})

}

func TestKnightCheck(t *testing.T) {
	var b board.Board
	expected := true
	b.ImportFEN("8/8/4k3/2N5/8/3K4/8/8 b - - 0 1")

	t.Run("Knight Check", func(t *testing.T) {
		if b.IsInCheck(board.BlackToMove) != expected {
			t.Errorf("Check not detected")
		}
	})

	b.ImportFEN("8/8/4k3/8/2N5/3K4/8/8 b - - 0 1")
	t.Run("No Knight Check", func(t *testing.T) {
		if b.IsInCheck(board.BlackToMove) == expected {
			t.Errorf("False Check detected")
		}
	})

}

func TestPawnCheck(t *testing.T) {
	var b board.Board
	expected := true
	b.ImportFEN("8/8/8/3k4/2P5/3K4/8/8 b - - 0 1")

	t.Run("Pawn Check", func(t *testing.T) {
		if b.IsInCheck(board.BlackToMove) != expected {
			t.Errorf("Check not detected")
		}
	})

	b.ImportFEN("88/8/8/3k4/4p3/3K4/8/8 w - - 0 1")
	t.Run("Black pawn Check", func(t *testing.T) {
		if b.IsInCheck(board.WhiteToMove) != expected {
			t.Errorf("False Check detected")
		}
	})

	b.ImportFEN("8/8/8/3k4/3P4/3K4/8/8 b - - 0 1")

	t.Run("Pawn Check", func(t *testing.T) {
		if b.IsInCheck(board.BlackToMove) == expected {
			t.Errorf("False check detected")
		}
	})

}
