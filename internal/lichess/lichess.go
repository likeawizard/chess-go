package lichess

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/likeawizard/chess-go/internal/board"
	eval "github.com/likeawizard/chess-go/internal/evaluation"
)

const (
	url = "https://lichess.org/api"

	accountPath   = "/account"
	challengePath = "/challenge"
	gamesPath     = "/playing"
)

type LichessConnector struct {
	Client *http.Client
	token  string
}

func (lc *LichessConnector) request(path, method string) ([]byte, error) {
	req, err := http.NewRequest(method, url+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+lc.token)
	resp, err := lc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else {
		return body, nil
	}
}

func NewLichessConnector() *LichessConnector {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env: %s\n", err)
		return nil
	}
	token := os.Getenv("LICHESS_TOKEN")

	return &LichessConnector{
		Client: &http.Client{},
		token:  token,
	}
}

func (lc *LichessConnector) CheckActiveGames() []NowPlaying {
	method := http.MethodGet
	path := fmt.Sprintf("%s%s", accountPath, gamesPath)
	body, err := lc.request(path, method)
	if err != nil {
		return nil
	}

	games := &GamesRsp{}
	err = json.Unmarshal(body, &games)
	if err != nil {
		return nil
	} else {
		return games.NowPlaying
	}
}

func (lc *LichessConnector) HandleActiveGames(games []NowPlaying) {
	for _, game := range games {
		if !game.IsMyTurn {
			continue
		}

		b := &board.Board{}
		b.Init()
		b.ImportFEN(game.Fen)
		e, _ := eval.NewEvalEngine(b)
		e.GetMove()
		best := e.RootNode.PickBestMove(b.SideToMove)
		move := best.MoveToPlay
		lc.MakeMove(game.GameID, move)
	}
}

func (lc *LichessConnector) MakeMove(gameID, move string) error {
	method := http.MethodPost
	path := fmt.Sprintf("/bot/game/%s/move/%s", gameID, move)

	_, err := lc.request(path, method)

	return err
}

func (lc *LichessConnector) ResignGame(gameID string) error {
	method := http.MethodPost
	path := fmt.Sprintf("/bot/game/%s/resign", gameID)

	_, err := lc.request(path, method)

	return err
}

func (lc *LichessConnector) GetChallenges() ([]Challenge, error) {
	method := http.MethodGet
	path := challengePath

	body, err := lc.request(path, method)
	if err != nil {
		return nil, err
	}

	challenges := &ChallengeRsp{}

	err = json.Unmarshal(body, &challenges)
	if err != nil {
		return nil, err
	} else {
		return challenges.In, nil
	}
}

func (lc *LichessConnector) ShouldAccept(ch Challenge) bool {
	return ch.TimeControl.Type == "unlimited" && ch.Variant.Key == "standard"
}

func (lc *LichessConnector) HandleChallenges(ch []Challenge) {
	for _, c := range ch {
		if lc.ShouldAccept(c) {
			lc.Accept(c)
		} else {
			lc.Decline(c)
		}
	}
}

func (lc *LichessConnector) Accept(c Challenge) error {
	method := http.MethodPost
	path := fmt.Sprintf("%s/%s/accept", challengePath, c.ID)
	_, err := lc.request(path, method)
	return err

}

func (lc *LichessConnector) Decline(c Challenge) error {
	method := http.MethodPost
	path := fmt.Sprintf("%s/%s/decline", challengePath, c.ID)
	_, err := lc.request(path, method)
	return err
}
