package main

import (
	"fmt"
	"time"

	"github.com/likeawizard/chess-go/internal/board"
)

type perftTest struct {
	name  string
	depth int
	count int
}

var containsError bool

func main() {
	depth := 4
	start := time.Now()
	// board.PerftDebug("r3k2r/Pppp1ppp/1b3nbN/nPB5/B1P1P3/q4N2/Pp1P2PP/R2Q1RK1 b kq - 1 1", depth)
	perft1(depth)
	perft2(depth)
	perft3(depth)
	perft4(depth)
	perft5(depth)
	perft6(depth)

	fmt.Printf("\nRun time: %v\n", time.Since(start))
	if !containsError {
		fmt.Println("All tests passed successfully.")
	} else {
		fmt.Println("Encountered errors")
	}
}

func test(fen string, depth int, data []perftTest) {
	for _, perft := range data {
		if depth == 0 {
			return
		}
		leafs, perf := board.Perft(fen, perft.depth)
		fmt.Printf("%s: %d %v %v\n", perft.name, leafs, perf, leafs == perft.count)
		if leafs != perft.count {
			containsError = true
		}
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
		{name: "depth8", depth: 8, count: 84998978956},
		{name: "depth9", depth: 9, count: 2439530234167},
		{name: "depth10", depth: 10, count: 69352859712417},
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
		{name: "depth7", depth: 7, count: 178633661},
		{name: "depth8", depth: 8, count: 3009794393},
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
	fmt.Println("Test 4")
	test(fen, depth, tests)
}

func perft5(depth int) {
	fen := "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"
	tests := []perftTest{
		{name: "depth1", depth: 1, count: 44},
		{name: "depth2", depth: 2, count: 1486},
		{name: "depth3", depth: 3, count: 62379},
		{name: "depth4", depth: 4, count: 2103487},
		{name: "depth5", depth: 5, count: 89941194},
	}
	fmt.Println("Test 5")
	test(fen, depth, tests)
}

func perft6(depth int) {
	fen := "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"
	tests := []perftTest{
		{name: "depth1", depth: 1, count: 46},
		{name: "depth2", depth: 2, count: 2079},
		{name: "depth3", depth: 3, count: 89890},
		{name: "depth4", depth: 4, count: 3894594},
		{name: "depth5", depth: 5, count: 164075551},
		{name: "depth6", depth: 6, count: 6923051137},
		{name: "depth7", depth: 7, count: 287188994746},
		{name: "depth8", depth: 8, count: 11923589843526},
		{name: "depth9", depth: 9, count: 490154852788714},
	}
	fmt.Println("Test 6")
	test(fen, depth, tests)
}
