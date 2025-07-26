package memory

import (
	"context"
	"github.com/google/uuid"
	"t03/internal/domain"
	"time"
)

type GameRepositoryImpl struct {
	storage *Storage
}

func NewGameRepository(storage *Storage) domain.GameRepository {
	return &GameRepositoryImpl{storage: storage}
}

const saveGameQuery = `
	INSERT INTO game_sessions (id, board_state, mode, player_x, player_o, state, turn, winner)
	VALUES ($1, $2, $3, $4, $5, $6, $7,$8)
	ON CONFLICT (id) DO UPDATE
	SET board_state = EXCLUDED.board_state,
	    player_o = EXCLUDED.player_o,
	    state     = EXCLUDED.state,
	    turn      = EXCLUDED.turn,
	    winner    = EXCLUDED.winner
`

const getGameQuery = `
		SELECT id, board_state, mode, player_x, player_o, state, turn,  winner
		FROM game_sessions
		WHERE id = $1
	`
const getAvalableGamesQuery = `
    SELECT id
    FROM game_sessions
    WHERE state IN (0,1)
      AND (
           player_x = $1
        OR player_o = $1
        OR (player_o = $2
            AND player_x <> $1 AND mode<>1)
      )`

const statsQuery = `
WITH my_games AS (
    SELECT
        id,
        CASE
            WHEN player_x = $1 THEN 'X'
            WHEN player_o = $1 THEN 'O'
        END AS role,
        winner,
        state 
    FROM
        game_sessions
    WHERE
        player_x = $1
        OR player_o = $1
)
SELECT
    COUNT(*) AS total_games,
    SUM(
        CASE
            WHEN winner = $1 THEN 1
            ELSE 0
        END
    ) AS wins,
    SUM(
        CASE
            WHEN winner <> $1
            AND state = 3 THEN 1
            ELSE 0
        END
    ) AS losses,
    SUM(
        CASE
            WHEN state = 2 THEN 1
            ELSE 0
        END
    ) AS draws,
    ROUND(
        100.0 * SUM(
            CASE
                WHEN winner = $1 THEN 1
                ELSE 0
            END
        ) / NULLIF(COUNT(*), 0),
        1
    ) AS win_rate_pct
FROM
    my_games;
`

func (repo *GameRepositoryImpl) GetPlayerStats(playerID uuid.UUID) (*domain.Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	var s domain.Stats
	if err := repo.storage.pool.QueryRow(ctx, statsQuery, playerID).Scan(
		&s.TotalGames, &s.Wins, &s.Losses, &s.Draws, &s.WinRatePct,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func (repo *GameRepositoryImpl) SaveGame(game *domain.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	entity := toEntity(game)

	_, err := repo.storage.pool.Exec(ctx, saveGameQuery, entity.GameId, entity.Board, entity.Mode, entity.Player_X, entity.Player_O, entity.State, entity.CurrentPID, entity.WinnerPID)

	return err
}

func (repo *GameRepositoryImpl) GetGame(id string) (*domain.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var entity GameEntity

	err := repo.storage.pool.QueryRow(ctx, getGameQuery, id).Scan(&entity.GameId, &entity.Board, &entity.Mode, &entity.Player_X, &entity.Player_O, &entity.State, &entity.CurrentPID, &entity.WinnerPID)

	if err != nil {
		return nil, err
	}

	return toDomain(&entity)
}
func (repo *GameRepositoryImpl) GetAvailableGames(pid string) (*domain.GamesList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	zero := uuid.Nil

	rows, err := repo.storage.pool.Query(ctx, getAvalableGamesQuery, pid, zero)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids uuid.UUIDs
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ToDomainGamesList(ids), nil
}
