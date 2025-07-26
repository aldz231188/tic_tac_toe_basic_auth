package dto

type GameRequest struct {
	Board [][]string `json:"board"`
	Mode  string     `json:"mode"`
}

type GameResponse struct {
	GameId    string     `json:"id"`
	Board     [][]string `json:"board"`
	PlayerXId string     `json:"playerX"`
	PlayerOId string     `json:"playerO"`
	Status    string     `json:"message"`
}

type Stats struct {
	TotalGames int     `json:"totalGames"`
	Wins       int     `json:"wins"`
	Losses     int     `json:"losses"`
	Draws      int     `json:"draws"`
	WinRatePct float64 `json:"winrate"`
}
