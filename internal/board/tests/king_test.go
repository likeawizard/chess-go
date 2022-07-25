package board

import (
	"reflect"
	"sort"
	"testing"

	"github.com/likeawizard/chess-go/internal/board"
)

//rnbqkb1r/pppp1ppp/5n2/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR w KQkq - 2 3

func TestSafeKingMoves(t *testing.T) {
	var b board.Board
	mExpected := []string{"e1e2", "e1f1"}
	b.ImportFEN("rnbqkb1r/pppp1ppp/5n2/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR w KQkq - 2 3")
	e1 := board.Coord{File: 4, Rank: 0}
	m, c := b.GetAvailableMoves(e1)

	t.Run("Two King moves", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}
	})

	t.Run("No King captures in position", func(t *testing.T) {
		if len(c) > 0 {
			t.Errorf("Got %v Expected empty array", c)
		}
	})
}

func TestKingCastlingMoves(t *testing.T) {
	var b board.Board
	mExpected := []string{"e1e2", "e1f1", "e1d1", "e1c1", "e1g1"}
	b.ImportFEN("r3k2r/pppqbppp/2np1n2/4p3/2BPP1b1/2NQ1N2/PPPB1PPP/R3K2R w KQkq - 5 8")
	e1 := board.Coord{File: 4, Rank: 0}
	m, c := b.GetAvailableMoves(e1)

	t.Run("Three regular King moves and both castling", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}
	})

	t.Run("No King captures in position", func(t *testing.T) {
		if len(c) > 0 {
			t.Errorf("Got %v Expected empty array", c)
		}
	})
}

func TestKingCastlingDenied(t *testing.T) {
	var b board.Board
	mExpected := []string{"e1e2", "e1f1", "e1d1"} // e1d1 is illegal
	b.ImportFEN("r3k2r/pppqbppp/2np1n2/4p3/2BPP3/2NQ1b2/PPPB1PPP/R3K1R1 w Qkq - 0 9")
	e1 := board.Coord{File: 4, Rank: 0}
	m, c := b.GetAvailableMoves(e1)

	t.Run("Two king moves no castling", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}
	})

	t.Run("No King captures in position", func(t *testing.T) {
		if len(c) > 0 {
			t.Errorf("Got %v Expected empty array", c)
		}
	})
}

func TestKingInCheckCantCastle(t *testing.T) {
	var b board.Board
	mExpected := []string{"e1f1", "e1d2"} // e1d2 is illegal
	b.ImportFEN("rnbqk1nr/pppp1ppp/8/4p3/1b2P3/3P1N2/PPP1BPPP/RNBQK2R w KQkq - 1 5")
	e1 := board.Coord{File: 4, Rank: 0}
	m, c := b.GetAvailableMoves(e1)

	t.Run("Cant castle under check", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}
	})

	t.Run("No King captures in position", func(t *testing.T) {
		if len(c) > 0 {
			t.Errorf("Got %v Expected empty array", c)
		}
	})
}

func TestKingCastledAndPiecePositions(t *testing.T) {
	var b board.Board
	mExpected := []string{"e1f1", "e1g1"}
	b.ImportFEN("rnbqk1nr/pppp1ppp/8/2b1p3/8/4PN2/PPPPBPPP/RNBQK2R w KQkq - 4 4")
	e1 := board.Coord{File: 4, Rank: 0}
	m, c := b.GetAvailableMoves(e1)
	g1 := board.Coord{File: 6, Rank: 0}
	f1 := board.Coord{File: 5, Rank: 0}

	t.Run("Verify castling", func(t *testing.T) {
		if !b.IsCastling("e1g1") {
			t.Errorf("Move not recognized as castling")
		}
	})

	b.MoveLongAlg("e1g1")

	t.Run("Safe castling", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}

		if b.AccessCoord(g1) != 6 || b.AccessCoord(f1) != 4 {
			t.Errorf("Expected 1 on g1 got: %d, Expected 4 on f1 got: %d", b.AccessCoord(g1), b.AccessCoord(f1))
		}

	})

	t.Run("No King captures in position", func(t *testing.T) {
		if len(c) > 0 {
			t.Errorf("Got %v Expected empty array", c)
		}
	})
}
