package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewPostgresPool(connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New((context.Background()), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
