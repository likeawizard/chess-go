package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/likeawizard/chess-go/internal/board"
	"github.com/likeawizard/chess-go/internal/render"

	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
		return
	}
	b1 := &board.Board{}
	b1.Init()
	depth, _ := strconv.Atoi(os.Getenv("EVALUATION_DEPTH"))
	b1.SetTrackMoves(true)
	board.InitEvalEngine(b1)
	var r render.BoardRender
	var aimove string
	var elapsed time.Duration

	r = render.New()
	r.InitRender(b1, &elapsed)
	go func() {
		for {
			start := time.Now()
			aimove = b1.GetMove(depth)
			elapsed = time.Since(start)
			b1.MoveLongAlg(aimove)
			r.Update()
			board.Evaluations = 0
			board.CachedEvals = 0
		}
	}()
	r.Run()
}
