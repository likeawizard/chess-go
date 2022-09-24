package board

// N, S, W, E, NW, NE, SW, SE
// 0:4 ranks & files, 4:8 diagonals
var Compass = []Square{8, -8, -1, 1, 7, 9, -9, -7}

// Number of squares to the edge in compass direction
var CompassBlock = [][]Square{}

var knightMoves = [][]Square{}

// pre calculate distances in all compass directions and possible knight jumps for every square
func init() {
	CompassBlock = make([][]Square, 64)
	min := func(a, b Square) Square {
		if a < b {
			return a
		}
		return b
	}
	for i := Square(0); i < 64; i++ {
		f, r := i%8, i/8
		n := 7 - r
		s := r
		w := f
		e := 7 - f
		CompassBlock[i] = []Square{n, s, w, e, min(n, w), min(n, e), min(s, w), min(s, e)}
	}
	preCalculateKnightMoves()
}

func preCalculateKnightMoves() {
	knightMoves = make([][]Square, 64)
	for c, comp := range CompassBlock {
		moves := make([]Square, 0)
		//2NW
		if comp[0] > 1 && comp[2] > 0 {
			moves = append(moves, Square(c)+15)
		}
		//2NE
		if comp[0] > 1 && comp[3] > 0 {
			moves = append(moves, Square(c)+17)
		}
		//2SW
		if comp[1] > 1 && comp[2] > 0 {
			moves = append(moves, Square(c)-17)
		}
		//2SE
		if comp[1] > 1 && comp[3] > 0 {
			moves = append(moves, Square(c)-15)
		}
		//2WN
		if comp[2] > 1 && comp[0] > 0 {
			moves = append(moves, Square(c)+6)
		}
		//2WS
		if comp[2] > 1 && comp[1] > 0 {
			moves = append(moves, Square(c)-10)
		}
		//2EN
		if comp[3] > 1 && comp[0] > 0 {
			moves = append(moves, Square(c)+10)
		}
		//2ES
		if comp[3] > 1 && comp[1] > 0 {
			moves = append(moves, Square(c)-6)
		}
		knightMoves[c] = moves
	}
}
