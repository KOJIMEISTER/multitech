package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidData     = errors.New("Invalid data")
	ErrSessionExists   = errors.New("Session already exists")
	ErrSessionNotFound = errors.New("nvalid or expired session")
)

type SessionsRepository interface {
	StoreSession(ctx context.Context, token string, userID uint, ttl time.Duration) error
	GetSession(ctx context.Context, token string) (uint, error)
	DeleteSession(ctx context.Context, token string) error
}
