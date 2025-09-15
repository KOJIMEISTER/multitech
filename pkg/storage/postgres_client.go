package storage

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPostgres() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	return pool, err
}
