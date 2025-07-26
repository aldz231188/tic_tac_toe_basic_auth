package memory

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"t03/internal/domain"
)

func toEntity(game *domain.Game) *GameEntity {
	var boardBuilder strings.Builder
	for i := range game.Board {
		for j := range game.Board[i] {
			boardBuilder.WriteByte('0' + byte(game.Board[i][j]))
		}
	}

	return &GameEntity{
		GameId:     game.GameId,
		Board:      boardBuilder.String(),
		Mode:       int(game.Mode),
		Player_X:   game.Player_X,
		Player_O:   game.Player_O,
		State:      int(game.State),
		CurrentPID: game.CurrentPID,
		WinnerPID:  game.WinnerPID,
	}
}

func toDomain(entity *GameEntity) (*domain.Game, error) {
	k := 0
	board := domain.Board{}

	for i := range board {
		for j := range board[i] {
			if k >= len(entity.Board) {
				return nil, errors.New("board data is too short")
			}
			c := entity.Board[k]
			if c < '0' || c > '9' {
				return nil, errors.New("invalid board character")
			}
			board[i][j] = domain.Cell(c - '0')
			k++
		}
	}

	return &domain.Game{
		GameId:     entity.GameId,
		Board:      board,
		Mode:       domain.Gametype(entity.Mode),
		Player_X:   entity.Player_X,
		Player_O:   entity.Player_O,
		State:      domain.GameState(entity.State),
		CurrentPID: entity.CurrentPID,
		WinnerPID:  entity.WinnerPID,
	}, nil
}

func ToDomainGamesList(games uuid.UUIDs) *domain.GamesList {
	return &domain.GamesList{
		Games: games,
	}

}
