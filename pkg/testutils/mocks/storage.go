package mocks

import (
	"context"
	"multitech/internal/models"
	"time"
)

const (
	ErrTestUserNotFound    = "test user not found"
	ErrTestSessionConflict = "test session conflict"
)

type MockUserRepository struct {
	GetUserByUsernameFunc func(ctx context.Context, username string) (*models.User, error)
	StoreUserFunc         func(ctx context.Context, user *models.User) error
}

func (mock *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return mock.GetUserByUsernameFunc(ctx, username)
}

func (mock *MockUserRepository) StoreUser(ctx context.Context, user *models.User) error {
	return mock.StoreUserFunc(ctx, user)
}

type MockSessionsRepository struct {
	StoreSessionFunc  func(ctx context.Context, token string, userID uint, duration time.Duration) error
	GetSessionFunc    func(ctx context.Context, token string) (uint, error)
	DeleteSessionFunc func(ctx context.Context, token string) error
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

func NewDefaultUserMock() *MockUserRepository {
	return &MockUserRepository{
		GetUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			return &models.User{
				ID:       1,
				Username: username,
				Password: "testpass",
			}, nil
		},
	}
}
