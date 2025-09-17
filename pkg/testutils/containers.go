package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestContainers struct {
	PostgresContainer *postgres.PostgresContainer
	RedisContainer    *redis.RedisContainer
	PostgresDSN       string
	RedisURL          string
}

func SetupContainers(ctx context.Context) (*TestContainers, error) {
	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(30*time.Second),
		))
	if err != nil {
		return nil, fmt.Errorf("Failed to start postgres: %w", err)
	}

	pgDSN, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get postgres DSN: %w", err)
	}

	redisContainer, err := redis.Run(ctx, "redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").WithStartupTimeout(30*time.Second),
		))
	if err != nil {
		return nil, fmt.Errorf("Failed to start redis: %w", err)
	}

	redisUrl, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get redis URL: %w", err)
	}

	return &TestContainers{
		PostgresContainer: pgContainer,
		RedisContainer:    redisContainer,
		PostgresDSN:       pgDSN,
		RedisURL:          redisUrl,
	}, nil
}

func (tc *TestContainers) Terminate(ctx context.Context) error {
	var errs []error

	if tc.PostgresContainer != nil {
		if err := tc.PostgresContainer.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("Failed to terminate Postgres: %w", err))
		}
	}

	if tc.RedisContainer != nil {
		if err := tc.RedisContainer.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("Failed to terminate Redis: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Container termination errors: %v", errs)
	}
	return nil
}
