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
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*15*1000)
			// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500*1000)
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
			moves = append(moves, best.MoveToPlay)

			b1.WritePGNToFile(b1.GeneratePGN(moves), "./dump.pgn")
			r.Update(best.MoveToPlay)
		}
	}()
	r.Run()
}
