package mocks

import (
	"context"
	"time"
)

type MockSessionsRepository struct {
	StoreSessionFunc  func(ctx context.Context, token string, userID uint, duration time.Duration) error
	GetSessionFunc    func(ctx context.Context, token string) (uint, error)
	DeleteSessionFunc func(ctx context.Context, token string) error
}

func NewDefaultSessionsMock() *MockSessionsRepository {
	return &MockSessionsRepository{
		StoreSessionFunc: func(ctx context.Context, token string, userID uint, duration time.Duration) error {
			return nil
		},
		GetSessionFunc: func(ctx context.Context, token string) (uint, error) {
			return 1, nil
		},
		DeleteSessionFunc: func(ctx context.Context, token string) error {
			return nil
		},
	}
}

func (mock *MockSessionsRepository) StoreSession(ctx context.Context, token string, userID uint, duration time.Duration) error {
	return mock.StoreSessionFunc(ctx, token, userID, duration)
}
func (mock *MockSessionsRepository) GetSession(ctx context.Context, token string) (uint, error) {
	return mock.GetSessionFunc(ctx, token)
}
func (mock *MockSessionsRepository) DeleteSession(ctx context.Context, token string) error {
	return mock.DeleteSessionFunc(ctx, token)
}
