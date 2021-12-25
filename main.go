package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)
import "Chess/board"

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
		return
	}
	var b1 *board.Board
	b1 = &board.Board{}
	b1.Init()
	depth, _ := strconv.Atoi(os.Getenv("EVALUATION_DEPTH"))
	b1.SetTrackMoves(true)
	board.InitEvalEngine(b1)
	var render *board.BoardRender
	var aimove string
	var elapsed time.Duration

	render = board.New().InitRender(b1, &elapsed)
	go func() {
		for {
			start := time.Now()
			aimove = b1.GetMove(depth)
			elapsed = time.Since(start)
			b1.MoveLongAlg(aimove)
			render.Update()
			board.Evaluations = 0
			board.CachedEvals = 0
		}
	}()
	render.Run()

	//fmt.Scanln()
}
