package board

import (
	"reflect"
	"sort"
	"testing"

	"github.com/likeawizard/chess-go/internal/board"
)

func TestRook(t *testing.T) {
	var b board.Board
	mExpected := []string{"e1e2", "e1e3", "e1e4", "e1e5", "e1f1"}
	cExpected := []string{"e1e6"}
	b.ImportFEN("3q3r/1r1P1p1k/1p1Qn1p1/1p6/1Pp2Pp1/7P/P1P5/3RR1K1 w - - 0 27")
	e1 := board.Coord{File: 4, Rank: 0}
	m, c := b.GetAvailableMoves(e1)

	t.Run("Rook moves across files / ranks", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}
	})

	t.Run("Rook captures Knight on e6", func(t *testing.T) {
		sort.Strings(c)
		sort.Strings(cExpected)
		if !reflect.DeepEqual(c, cExpected) {
			t.Errorf("Got %v Expected %v", c, cExpected)
		}
	})
}
