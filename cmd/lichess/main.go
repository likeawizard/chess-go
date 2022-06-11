package main

import (
	"fmt"
	"time"

	"github.com/likeawizard/chess-go/internal/lichess"
)

func main() {
	lc := lichess.NewLichessConnector()
	for {
		challenges, err := lc.GetChallenges()
		if err != nil {
			fmt.Println(err)
		}

		if len(challenges) > 0 {
			fmt.Println(challenges)
			lc.HandleChallenges(challenges)
		}

		games := lc.CheckActiveGames()

		if len(games) > 0 {
			fmt.Println(games)
			lc.HandleActiveGames(games)
		}
		time.Sleep(5 * time.Second)
	}
}
