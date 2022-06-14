package main

import (
	"fmt"

	"github.com/likeawizard/chess-go/internal/lichess"
)

func main() {
	lc := lichess.NewLichessConnector()
	decoder, err := lc.OpenEventStream()
	if err != nil {
		fmt.Printf("Failed to open stream: %s\n", err)
		return
	}

	go lc.ListenToEvents(decoder)
	lc.HandleMoveQueue()
}
