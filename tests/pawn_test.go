package main

import (
	"github.com/likeawizard/chess-go/board"
	"reflect"
	"sort"
	"testing"
)

func TestFirstMove(t *testing.T) {
	var b board.Board
	mExpected := []string{"e2e3", "e2e4"}
	b.Init()
	e2 := board.Coord{File: 4, Rank: 1}
	m, c := b.GetAvailableMoves(e2)

	t.Run("First move moves", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}
	})

	t.Run("First move captures", func(t *testing.T) {
		if len(c) > 0 {
			t.Errorf("Got %v Expected empty array", c)
		}
	})
}

func TestAlreadyMoved(t *testing.T) {
	var b board.Board
	position := "rnbqkbnr/ppp1pppp/8/3p4/8/4P3/PPPP1PPP/RNBQKBNR w KQkq - 0 2"
	mExpected := []string{"e3e4"}
	b.ImportFEN(position)
	e3 := board.Coord{File: 4, Rank: 2}

	m, c := b.GetAvailableMoves(e3)

	t.Run("Pawn already moved moves", func(t *testing.T) {
		sort.Strings(m)
		sort.Strings(mExpected)
		if !reflect.DeepEqual(m, mExpected) {
			t.Errorf("Got %v Expected %v", m, mExpected)
		}
	})

	t.Run("Pawn already moved captures", func(t *testing.T) {
		if len(c) > 0 {
			t.Errorf("Got %v Expected empty array", c)
		}
	})
}

func TestBlockedPawnTwoCaptures(t *testing.T) {
	var b board.Board
	position := "rnbqkbnr/ppp3pp/8/3ppp2/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 4"
	cExpected := []string{"e4d5", "e4f5"}
	b.ImportFEN(position)
	e3 := board.Coord{File: 4, Rank: 3}

	m, c := b.GetAvailableMoves(e3)

	t.Run("Two pawn captures", func(t *testing.T) {
		sort.Strings(c)
		sort.Strings(cExpected)
		if !reflect.DeepEqual(c, cExpected) {
			t.Errorf("Got %v Expected %v", c, cExpected)
		}
	})

	t.Run("Pawn move forward blocked", func(t *testing.T) {
		if len(m) > 0 {
			t.Errorf("Got %v len %d Expected empty array", m, len(m))
		}
	})
}

func TestEnPassant(t *testing.T) {
	var b board.Board
	position := "rnbqkbnr/ppp2ppp/4p3/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3"
	cExpected := []string{"e5d6"}
	b.ImportFEN(position)
	e5 := board.Coord{File: 4, Rank: 4}
	epTarget := board.Coord{File: 3, Rank: 4}

	_, c := b.GetAvailableMoves(e5)
	b.MoveLongAlg("e5d6")

	t.Run("En Passant Capture", func(t *testing.T) {
		if !reflect.DeepEqual(c, cExpected) {
			t.Errorf("Got %v Expected %v", c, cExpected)
		}
		if b.AccessCoord(epTarget) != 0 {
			t.Errorf("Pawn not removed after capture got: %v", b.AccessCoord(epTarget))
		}
	})
}
