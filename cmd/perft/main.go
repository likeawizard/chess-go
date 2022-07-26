package main

import (
	"fmt"

	"github.com/likeawizard/chess-go/internal/board"
)

type perftTest struct {
	name  string
	depth int
	count int
}

func main() {
	depth := 6
	// board.PerftDebug("r3k2N/p1ppq1b1/1n2pnp1/1b1P4/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQq - 0 2", depth)
	// perft1(depth)
	// perft2(depth)
	// perft3(8)
	perft4(depth)
}

func test(fen string, depth int, data []perftTest) {
	for _, perft := range data {
		if depth == 0 {
			return
		}
		leafs, perf := board.Perft(fen, perft.depth)
		fmt.Printf("%s: %d %v %v\n", perft.name, leafs, perf, leafs == perft.count)
		depth--
	}
	fmt.Println()
}

func perft1(depth int) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	tests := []perftTest{
		{name: "depth1", depth: 1, count: 20},
		{name: "depth2", depth: 2, count: 400},
		{name: "depth3", depth: 3, count: 8902},
		{name: "depth4", depth: 4, count: 197281},
		{name: "depth5", depth: 5, count: 4865609},
		{name: "depth6", depth: 6, count: 119060324},
		{name: "depth7", depth: 7, count: 3195901860},
	}
	fmt.Println("Test 1")
	test(fen, depth, tests)
}

func perft2(depth int) {
	fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
	tests := []perftTest{
		{name: "depth1", depth: 1, count: 48},
		{name: "depth2", depth: 2, count: 2039},
		{name: "depth3", depth: 3, count: 97862},
		{name: "depth4", depth: 4, count: 4085603},
		{name: "depth5", depth: 5, count: 193690690},
		{name: "depth6", depth: 6, count: 8031647685},
	}
	fmt.Println("Test 2")
	test(fen, depth, tests)
}

func perft3(depth int) {
	fen := "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1"
	tests := []perftTest{
		{name: "depth1", depth: 1, count: 14},
		{name: "depth2", depth: 2, count: 191},
		{name: "depth3", depth: 3, count: 2812},
		{name: "depth4", depth: 4, count: 43238},
		{name: "depth5", depth: 5, count: 674624},
		{name: "depth6", depth: 6, count: 11030083},
		{name: "depth6", depth: 7, count: 178633661},
		{name: "depth6", depth: 8, count: 3009794393},
	}
	fmt.Println("Test 3")
	test(fen, depth, tests)
}

func perft4(depth int) {
	fen := "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"
	tests := []perftTest{
		{name: "depth1", depth: 1, count: 6},
		{name: "depth2", depth: 2, count: 264},
		{name: "depth3", depth: 3, count: 9467},
		{name: "depth4", depth: 4, count: 422333},
		{name: "depth5", depth: 5, count: 15833292},
		{name: "depth6", depth: 6, count: 706045033},
	}
	fmt.Println("Test 3")
	test(fen, depth, tests)
}
