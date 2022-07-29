package main

import (
	"context"
	"fmt"
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
	moves := make([]board.Move, 0)
	r := render.New(cfg)
	r.InitRender(b1, e)
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10*1000)
			// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500*1000)
			move := e.GetMove(ctx)
			defer cancel()
			if move == 0 {
				fmt.Println("No legal moves.")
				return
			}
			b1.MoveLongAlg(move)
			moves = append(moves, move)

			b1.WritePGNToFile(b1.GeneratePGN(moves), "./dump.pgn")
			r.Update(move)
		}
	}()
	r.Run()
}
