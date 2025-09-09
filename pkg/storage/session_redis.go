package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) SessionsRepository {
	return &sessionRepository{
		client: client,
	}
}

func (sessRepo *sessionRepository) StoreSession(ctx context.Context, token string, userID uint, ttl time.Duration) error {
	key := sessionKey(token)
	exists, err := sessRepo.client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis error: %w", err)
	}
	if exists == 1 {
		return ErrSessionExists
	}

	return sessRepo.client.SetEx(ctx, key, userID, ttl).Err()
}

func (sessRepo *sessionRepository) GetSession(ctx context.Context, token string) (uint, error) {
	userID, err := sessRepo.client.Get(ctx, sessionKey(token)).Uint64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, ErrSessionNotFound
		}
		return 0, fmt.Errorf("redis error: %w", err)
	}
	return uint(userID), err
}

func (sessRepo *sessionRepository) DeleteSession(ctx context.Context, token string) error {
	return sessRepo.client.Del(ctx, sessionKey(token)).Err()
}

func sessionKey(token string) string {
	return "token:" + token
}
