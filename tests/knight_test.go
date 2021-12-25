package main

import (
	"Chess/board"
	"reflect"
	"sort"
	"testing"
)

func TestKnight(t *testing.T) {
	var b board.Board
	mExpected := []string{"e4d6", "e4f6", "e4c3", "e4g3"}
	cExpected := []string{"e4c5", "e4g5"}
	b.ImportFEN("rnbqkbnr/pp1ppp1p/8/2p3p1/4N3/8/PPPPPPPP/R1BQKBNR w KQkq - 0 3")
	e4 := board.Coord{File: 4, Rank: 3}
	m, c := b.GetAvailableMoves(e4)

	t.Run("Four Knight moves", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}
	})

	t.Run("Two knight captures", func(t *testing.T) {
		sort.Strings(c)
		sort.Strings(cExpected)
		if !reflect.DeepEqual(c, cExpected) {
			t.Errorf("Got %v Expected %v", c, cExpected)
		}
	})
}
