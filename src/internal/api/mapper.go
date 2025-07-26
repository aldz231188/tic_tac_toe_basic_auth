package api

import (
	"errors"
	"t03/internal/api/dto"
	"t03/internal/domain"

	"github.com/google/uuid"
)

func ToDomainGame(pasedId string, s dto.GameRequest) (*domain.Game, error) {
	var game domain.Game

	id, err := uuid.Parse(pasedId)
	if err != nil {
		return &game, err
	}
	game.GameId = id

	switch s.Mode {
	case "human":
		game.Mode = domain.PVP
	case "ai":
		game.Mode = domain.PVE
	}

	if len(s.Board) != 3 {
		return &game, errors.New("board must have 3 rows")
	}

	for i := 0; i < 3; i++ {
		if len(s.Board[i]) != 3 {
			return &game, errors.New("each row must have 3 columns")
		}

		for j := 0; j < 3; j++ {
			switch s.Board[i][j] {
			case "X":
				game.Board[i][j] = domain.X
			case "O":
				game.Board[i][j] = domain.O
			case "":
				game.Board[i][j] = domain.Empty
			}
		}
	}

	return &game, nil
}

func ToGameResponse(game *domain.Game) dto.GameResponse {
	board := make([][]string, 3)
	for i := range board {
		board[i] = make([]string, 3)
		for j := 0; j < 3; j++ {
			switch game.Board[i][j] {
			case domain.X:
				board[i][j] = "X"
			case domain.O:
				board[i][j] = "O"
			default:
				board[i][j] = ""
			}
		}
	}

	var message string
	switch game.State {
	case domain.StatusWaiting:
		message = "Waiting for another player to connect"
	case domain.StatusTurn:
		message = "In progress"
	case domain.StatusDraw:
		message = "Draw"
	case domain.StatusWin:
		message = "Player " + game.WinnerPID.String() + " won"
	}

	return dto.GameResponse{
		GameId:    game.GameId.String(),
		Board:     board,
		PlayerXId: game.Player_X.String(),
		PlayerOId: game.Player_O.String(),
		Status:    message,
	}
}

func ToGamesListResponse(games *domain.GamesList) []string {
	return games.Games.Strings()

}
func ToStats(stats *domain.Stats) *dto.Stats {
	return &dto.Stats{
		TotalGames: stats.TotalGames,
		Wins:       stats.Wins,
		Losses:     stats.Losses,
		Draws:      stats.Draws,
		WinRatePct: stats.WinRatePct,
	}

}
