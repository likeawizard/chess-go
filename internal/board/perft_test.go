package board

import "testing"

type perftTest struct {
	name  string
	depth int
	count int
}

var tt []perftTest = []perftTest{
	{name: "depth1", depth: 1, count: 20},
	{name: "depth2", depth: 2, count: 400},
	{name: "depth3", depth: 3, count: 8902},
	{name: "depth4", depth: 4, count: 197281},
	{name: "depth5", depth: 5, count: 4865609},
	// {name: "depth6", depth: 6, count: 119060324},
	// {name: "depth7", depth: 7, count: 3195901860},
}

func TestDepth1(t *testing.T) {
	for _, perft := range tt {
		t.Run(perft.name, func(t *testing.T) {
			leafs, _ := Perft(startingFEN, perft.depth)
			if leafs != perft.count {
				t.Fatalf("Perft failed: got - %v, expected %v", leafs, perft.count)
			}
		})
	}

}
