package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/likeawizard/chess-go/internal/board"
	eval "github.com/likeawizard/chess-go/internal/evaluation"
	"github.com/likeawizard/chess-go/internal/render"
	_ "go.uber.org/automaxprocs"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env: %s\n", err)
		return
	}
	b1 := &board.Board{}
	b1.Init()
	e, err := eval.NewEvalEngine(b1)
	if err != nil {
		fmt.Printf("Unable to load EvalEngine: %s\n", err)
		return
	}
	b1.SetTrackMoves(true)

	r := render.New()
	r.InitRender(b1, e)
	RegisterIterrupt(b1)
	go func() {
		for {
			e.GetMove()
			best := e.RootNode.PickBestMove(b1.SideToMove)
			candidates := e.RootNode.PickBestMoves(3)
			for _, move := range candidates {
				fmt.Printf("%.2f %v\n", move.Evaluation, move.ConstructLine())
			}
			b1.MoveLongAlg(best.MoveToPlay)
			e.PlayMove(best)
			r.Update()
		}
	}()
	r.Run()
}

func RegisterIterrupt(b *board.Board) {
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println(b.GetMoveList())
		os.Exit(0)
	}()
}
