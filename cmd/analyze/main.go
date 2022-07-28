package main

import (
	"context"
	"flag"
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
	fen := flag.String("fen", "", "FEN")
	flag.Parse()
	b := &board.Board{}
	b.ImportFEN(*fen)
	if b.ExportFEN() != *fen {
		fmt.Printf("Error importing FEN: %s, %s\n", b.ExportFEN(), *fen)
		return
	}
	e, err := eval.NewEvalEngine(b, cfg)
	if err != nil {
		fmt.Printf("Unable to load EvalEngine: %s\n", err)
		return
	}

	r := render.New(cfg)
	r.InitRender(b, e)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*15*1000)
	best := e.GetMove(ctx)
	defer cancel()
	candidates := e.RootNode.PickBestMoves(3)
	for _, move := range candidates {
		fmt.Printf("%.2f %v\n", move.Evaluation, move.ConstructLine())
	}
	b.MoveLongAlg(best.MoveToPlay)
	e.PlayMove(best)
	fmt.Println(best.MoveToPlay)
	r.Update(best.MoveToPlay)
}
