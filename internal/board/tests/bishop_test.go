package board

import (
	"reflect"
	"sort"
	"testing"

	"github.com/likeawizard/chess-go/internal/board"
)

func TestFirstMoveBlocked(t *testing.T) {
	var b board.Board
	b.Init()
	f1 := board.Coord{File: 5, Rank: 0}
	m, c := b.GetAvailableMoves(f1)

	t.Run("First move moves", func(t *testing.T) {
		if len(m) > 0 {
			t.Errorf("Got %v Expected empty array", m)
		}
	})

	t.Run("First move captures", func(t *testing.T) {
		if len(c) > 0 {
			t.Errorf("Got %v Expected empty array", c)
		}
	})
}

func TestOpenA6F1Diagonal(t *testing.T) {
	var b board.Board
	mExpected := []string{"f1e2", "f1d3", "f1c4", "f1b5"}
	cExpected := []string{"f1a6"}
	b.ImportFEN("rnbqkbnr/1ppppppp/p7/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2")
	f1 := board.Coord{File: 5, Rank: 0}
	m, c := b.GetAvailableMoves(f1)

	t.Run("Diagonal moves", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}
	})

	t.Run("Diagonal captures", func(t *testing.T) {
		sort.Strings(c)
		sort.Strings(cExpected)
		if !reflect.DeepEqual(c, cExpected) {
			t.Errorf("Got %v Expected %v", c, cExpected)
		}
	})
}
