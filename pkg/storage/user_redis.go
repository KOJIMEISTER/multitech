package storage

import (
	"context"
	"fmt"
	"multitech/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type userRepository struct {
	client *redis.Client
}

func NewUserRepository(client *redis.Client) UserRepository {
	return &userRepository{
		client: client,
	}
}

func (userRepo *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	exists, err := userRepo.client.Exists(ctx, userKey(user.Username)).Result()
	if err != nil {
		return fmt.Errorf("redis error: %w", err)
	}
	if exists == 1 {
		return ErrUserExists
	}

	userData := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"password":   user.Password,
		"created_at": user.CreatedAt.Format(time.RFC3339),
	}

	if err := userRepo.client.HSet(ctx, userKey(user.Username), userData).Err(); err != nil {
		return fmt.Errorf("redis error: %w", err)
	}

	return nil
}

func (userRepo *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	data, err := userRepo.client.HGetAll(ctx, userKey(username)).Result()
	if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}
	if len(data) == 0 {
		return nil, ErrUserNotFound
	}

	createdAt, err := time.Parse(time.RFC3339, data["created_at"])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid created_at format", ErrInvalidData)
	}

	return &models.User{
		ID:        uint(parseUint(data["id"])),
		Username:  data["username"],
		Email:     data["email"],
		Password:  data["password"],
		CreatedAt: createdAt,
	}, nil
}

func userKey(username string) string {
	return "user:" + username
}

func parseUint(str string) uint64 {
	var number uint64
	fmt.Sscanf(str, "%d", &number)
	return number
}
