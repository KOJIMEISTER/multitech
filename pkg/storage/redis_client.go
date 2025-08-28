package storage

import (
	"context"
	"errors"
	"fmt"
	"multitech/internal/models"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidData   = errors.New("invalid data")
	ErrSessionExists = errors.New("session already exists")
)

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	return err
}

func CreateUser(ctx context.Context, user *models.User) error {
	exists, err := RedisClient.Exists(ctx, userKey(user.Username)).Result()
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

	if err := RedisClient.HSet(ctx, userKey(user.Username), userData).Err(); err != nil {
		return fmt.Errorf("redis error: %w", err)
	}

	return nil
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	data, err := RedisClient.HGetAll(ctx, userKey(username)).Result()
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

func StoreSession(ctx context.Context, token string, userID uint, ttl time.Duration) error {
	key := sessionKey(token)
	exists, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis error: %w", err)
	}
	if exists == 1 {
		return ErrSessionExists
	}

	return RedisClient.SetEx(ctx, key, userID, ttl).Err()
}

func GetSession(ctx context.Context, token string) (uint, error) {
	userID, err := RedisClient.Get(ctx, sessionKey(token)).Uint64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("redis error: %w", err)
	}
	return uint(userID), err
}

func DeleteSession(ctx context.Context, token string) error {
	return RedisClient.Del(ctx, sessionKey(token)).Err()
}

func userKey(username string) string {
	return "user:" + username
}

func sessionKey(token string) string {
	return "token:" + token
}

func parseUint(str string) uint64 {
	var number uint64
	fmt.Sscanf(str, "%d", &number)
	return number
}
