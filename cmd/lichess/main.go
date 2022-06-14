package main

import (
	"fmt"
	"time"

	"github.com/likeawizard/chess-go/internal/lichess"
)

func main() {
	lc := lichess.NewLichessConnector()
	decoder, err := lc.OpenEventStream()
	if err != nil {
		fmt.Printf("Failed to open stream: %s\n", err)
		return
	}

	go lc.ListenToChallenges(decoder)
	for {
		games := lc.CheckActiveGames()

		if len(games) > 0 {
			fmt.Println(games)
			lc.HandleActiveGames(games)
		}
		time.Sleep(5 * time.Second)
	}
}
