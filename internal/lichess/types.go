package lichess

type ChallengeRsp struct {
	In  []Challenge `json:"in"`
	Out []Challenge `json:"out"`
}

type Challenge struct {
	ID          string      `json:"id"`
	URL         string      `json:"url"`
	Status      string      `json:"status"`
	Challenger  Challenger  `json:"challenger"`
	DestUser    DestUser    `json:"destUser"`
	Variant     Variant     `json:"variant"`
	Rated       bool        `json:"rated"`
	Speed       string      `json:"speed"`
	TimeControl TimeControl `json:"timeControl"`
	Color       string      `json:"color"`
	FinalColor  string      `json:"finalColor"`
	Perf        Perf        `json:"perf"`
	Direction   string      `json:"direction"`
}
type Challenger struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Title       interface{} `json:"title"`
	Rating      int         `json:"rating"`
	Provisional bool        `json:"provisional"`
	Online      bool        `json:"online"`
}
type DestUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Rating      int    `json:"rating"`
	Provisional bool   `json:"provisional"`
	Online      bool   `json:"online"`
}
type Variant struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Short string `json:"short"`
}
type TimeControl struct {
	Type string `json:"type"`
}
type Perf struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
}

type GamesRsp struct {
	NowPlaying []NowPlaying `json:"nowPlaying"`
}

// type Variant struct {
// 	Key string `json:"key"`
// 	Name string `json:"name"`
// }
type Opponent struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
}
type NowPlaying struct {
	FullID   string   `json:"fullId"`
	GameID   string   `json:"gameId"`
	Fen      string   `json:"fen"`
	Color    string   `json:"color"`
	LastMove string   `json:"lastMove"`
	Source   string   `json:"source"`
	Variant  Variant  `json:"variant"`
	Speed    string   `json:"speed"`
	Perf     string   `json:"perf"`
	Rated    bool     `json:"rated"`
	HasMoved bool     `json:"hasMoved"`
	Opponent Opponent `json:"opponent"`
	IsMyTurn bool     `json:"isMyTurn"`
}
