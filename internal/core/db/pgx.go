package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPgx(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect (pgx): %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("timescaledb unreachable: %w", err)
	}
	return pool, nil
}
