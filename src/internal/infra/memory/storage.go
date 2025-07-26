package memory

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Storage struct {
	pool *pgxpool.Pool
}

type Config struct {
	DSN string
}

func NewStorage(lc fx.Lifecycle, cfg Config) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return &Storage{pool: pool}, nil
}

func NewPGConfig() Config {
	return Config{
		DSN: "postgres://postgres:Qwaszx_1@localhost:5432/tic_tac_toe",
	}
}
