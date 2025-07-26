package domain

import (
	"github.com/google/uuid"
	"t03/internal/api/dto"
)

type GameService interface {
	PlayerVsAi(game *Game, playerId string) (*Game, error)
	PlayerMove(game *Game, playerId string) (*Game, error)
	NewGame(playerID string, gameType string) (string, error)
	GetAvailableGames(pid string) (*GamesList, error)
	ConnectToGame(gameId, userId string) (*Game, error)
	GetPlayerStats(playerID string) (*Stats, error)
}

type GameRepository interface {
	SaveGame(game *Game) error
	GetGame(id string) (*Game, error)
	GetAvailableGames(pid string) (*GamesList, error)
	SaveUser(user *User) error
	GetUser(login string) (*User, error)
	GetPlayerStats(playerID uuid.UUID) (*Stats, error)
}

type UserService interface {
	Register(request dto.SignUpRequest) (string, error)
	AuthenticateBasic(base64Credentials string) (string, error)
}
