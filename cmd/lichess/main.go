package main

import (
	"fmt"

	"github.com/likeawizard/chess-go/internal/config"
	"github.com/likeawizard/chess-go/internal/lichess"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load app config: %s\n", err)
	}
	lc := lichess.NewLichessConnector(cfg)
	decoder, err := lc.OpenEventStream()
	if err != nil {
		fmt.Printf("Failed to open stream: %s\n", err)
		return
	}

	lc.ListenToEvents(decoder)
}
