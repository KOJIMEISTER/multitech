package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidData   = errors.New("invalid data")
	ErrSessionExists = errors.New("session already exists")
)

type SessionsRepository interface {
	StoreSession(ctx context.Context, token string, userID uint, ttl time.Duration) error
	GetSession(ctx context.Context, token string) (uint, error)
	DeleteSession(ctx context.Context, token string) error
}
