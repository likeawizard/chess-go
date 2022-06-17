package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/likeawizard/chess-go/internal/board"
	"github.com/likeawizard/chess-go/internal/config"
	eval "github.com/likeawizard/chess-go/internal/evaluation"
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
	eval.DEBUG = true
	e, err := eval.NewEvalEngine(b, cfg)
	if err != nil {
		fmt.Printf("Unable to load EvalEngine: %s\n", err)
		return
	}

	fmt.Println(eval.GetEvaluation(e, b))
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
