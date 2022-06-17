package lichess

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/likeawizard/chess-go/internal/board"
	"github.com/likeawizard/chess-go/internal/config"
	eval "github.com/likeawizard/chess-go/internal/evaluation"
)

const (
	apiUrl = "https://lichess.org/api"

	accountPath   = "/account"
	challengePath = "/challenge"
	gamesPath     = "/playing"
)

type LichessConnector struct {
	Client    *http.Client
	token     string
	MoveQueue chan MoveQueue
	Config    *config.Config
}

func (lc *LichessConnector) request(path, method string, payload io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, apiUrl+path, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+lc.token)
	if payload != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
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

func NewLichessConnector(c *config.Config) *LichessConnector {
	return &LichessConnector{
		Client:    &http.Client{},
		token:     c.Lichess.APIToken,
		MoveQueue: make(chan MoveQueue),
		Config:    c,
	}
}

func (lc *LichessConnector) CheckActiveGames() []NowPlaying {
	method := http.MethodGet
	path := fmt.Sprintf("%s%s", accountPath, gamesPath)
	body, err := lc.request(path, method, nil)
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
		b.ImportFEN(game.Fen)
		e, _ := eval.NewEvalEngine(b, lc.Config)
		e.GetMove()
		best := e.RootNode.PickBestMove(b.SideToMove)
		move := best.MoveToPlay
		err := lc.MakeMove(game.GameID, move)
		if err != nil {
			lc.ResignGame(game.GameID)
		}
	}
}

func (lc *LichessConnector) MakeMove(gameID, move string) error {
	method := http.MethodPost
	path := fmt.Sprintf("/bot/game/%s/move/%s", gameID, move)

	_, err := lc.request(path, method, nil)

	return err
}

func (lc *LichessConnector) ResignGame(gameID string) error {
	method := http.MethodPost
	path := fmt.Sprintf("/bot/game/%s/resign", gameID)

	_, err := lc.request(path, method, nil)

	return err
}

func (lc *LichessConnector) GetChallenges() ([]Challenge, error) {
	method := http.MethodGet
	path := challengePath

	body, err := lc.request(path, method, nil)
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

func (lc *LichessConnector) ShouldAccept(ch Challenge) (bool, string) {
	if ch.Variant.Key != "standard" {
		return false, "variant"
	}
	if ch.TimeControl.Type != "unlimited" && ch.TimeControl.Type != "correspondence" {
		return false, "tooFast"
	}
	return true, ""
}

func (lc *LichessConnector) HandleChallenges(ch []Challenge) {
	for _, c := range ch {
		if ok, reason := lc.ShouldAccept(c); ok {
			lc.Accept(c)
		} else {
			lc.Decline(c, reason)
		}
	}
}

func (lc *LichessConnector) Accept(c Challenge) error {
	method := http.MethodPost
	path := fmt.Sprintf("%s/%s/accept", challengePath, c.ID)
	_, err := lc.request(path, method, nil)
	return err

}

func (lc *LichessConnector) Decline(c Challenge, reason string) error {
	method := http.MethodPost
	if reason == "" {
		reason = "generic"
	}
	data := url.Values{}
	data.Set("reason", reason)
	payload := data.Encode()

	path := fmt.Sprintf("%s/%s/decline", challengePath, c.ID)

	_, err := lc.request(path, method, strings.NewReader(payload))
	return err
}

func (lc *LichessConnector) OpenEventStream() (*json.Decoder, error) {
	req, err := http.NewRequest(http.MethodGet, apiUrl+"/stream/event", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+lc.token)
	resp, err := lc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder, nil
}

func (lc *LichessConnector) ListenToEvents(decoder *json.Decoder) {
	for decoder.More() {
		var e StreamEvent
		err := decoder.Decode(&e)
		if err != nil {
			fmt.Printf("Error decoding stream event: %s\n", err)
		}

		switch e.Type {
		case EVENT_CHALLENGE:
			fmt.Printf("New Challenge: %v\n", e.Challenge)
			lc.HandleChallenges([]Challenge{e.Challenge})
		case EVENT_GAME_START:
			fmt.Printf("New Game: %v\n", e.Game)
			go lc.ListenToGame(e.Game)
		default:
			fmt.Printf("Unhandled event: %s\n", e.Type)
		}
	}
}

func (lc *LichessConnector) ListenToGame(game Game) {
	//https://lichess.org/api/bot/game/stream/{gameId}
	path := fmt.Sprintf("%s/bot/game/stream/%s", apiUrl, game.GameID)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		fmt.Printf("Error opening GameStream: %s\n", err)
	}
	req.Header.Add("Authorization", "Bearer "+lc.token)
	resp, err := lc.Client.Do(req)
	if err != nil {
		fmt.Printf("Error opening GameStream: %s\n", err)
	}
	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var gs GameState
		err := decoder.Decode(&gs)
		if err != nil {
			fmt.Printf("Error decoding GameState: %s\n", err)
			continue
		}

		//TODO: Dirty stream implementation. Uses streams only as a notification to check active games.
		//Rewrite to use game event gameFull/gameState to independently construct game state.
		switch gs.Type {
		case GAME_EVENT_FULL, GAME_EVENT_STATE:
			g, err := lc.getActiveGame(game.GameID)
			if err != nil {
				fmt.Printf("Not my turn. Don't care....\n")
				continue
			}
			lc.MoveQueue <- MoveQueue{Fen: g.Fen, GameID: g.GameID}

		default:
			fmt.Printf("Unhandled game state: %s\n", gs.Type)
		}

	}
}

func (lc *LichessConnector) getActiveGame(gameID string) (*NowPlaying, error) {
	nowPlaying := lc.CheckActiveGames()
	for _, game := range nowPlaying {
		if game.GameID == gameID && game.IsMyTurn {
			return &game, nil
		}
	}
	return nil, fmt.Errorf("game not found")

}

func (lc *LichessConnector) HandleMoveQueue() {
	for g := range lc.MoveQueue {
		go func(g MoveQueue) {
			fmt.Printf("Pondering on move: %v\n", g)
			b := &board.Board{}
			b.ImportFEN(g.Fen)
			e, _ := eval.NewEvalEngine(b, lc.Config)
			e.GetMove()
			best := e.RootNode.PickBestMove(b.SideToMove)
			move := best.MoveToPlay

			fmt.Printf("Making %s move in %s (FEN: %s )\n", move, g.GameID, g.Fen)
			err := lc.MakeMove(g.GameID, move)
			if err != nil {
				fmt.Printf("Illegal move - resigning.\n")
				lc.ResignGame(g.GameID)
			}
		}(g)
	}
}
