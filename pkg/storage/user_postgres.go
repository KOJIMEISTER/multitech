package storage

import (
	"context"
	"multitech/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{pool: pool}
}

func (userRepo *postgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, nil
}

func (userRepo *postgresUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return nil
}
