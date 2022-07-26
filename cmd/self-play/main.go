package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/likeawizard/chess-go/internal/board"
	"github.com/likeawizard/chess-go/internal/config"
	eval "github.com/likeawizard/chess-go/internal/evaluation"
	"github.com/likeawizard/chess-go/internal/render"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load app config: %s\n", err)
	}
	b1 := &board.Board{}
	b1.Init(cfg)
	e, err := eval.NewEvalEngine(b1, cfg)
	if err != nil {
		fmt.Printf("Unable to load EvalEngine: %s\n", err)
		return
	}
	b1.SetTrackMoves(true)

	r := render.New(cfg)
	r.InitRender(b1, e)
	RegisterIterrupt(b1)
	go func() {
		for {
			// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*15*1000)
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500*1000)
			best := e.GetMove(ctx)
			defer cancel()
			// candidates := e.RootNode.PickBestMoves(3)
			// for _, move := range candidates {
			// 	fmt.Printf("%.2f %v\n", move.Evaluation, move.ConstructLine())
			// }
			if best == nil {
				fmt.Println("No legal moves.")
				return
			}
			b1.MoveLongAlg(best.MoveToPlay)
			e.PlayMove(best)
			b1.WritePGNToFile("./dump.pgn")
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
		fmt.Println(b.GeneratePGN())
		os.Exit(0)
	}()
}
