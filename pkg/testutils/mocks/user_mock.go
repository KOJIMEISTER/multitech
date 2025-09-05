package mocks

import (
	"context"
	"multitech/internal/models"
)

type MockUserRepository struct {
	GetUserByUsernameFunc func(ctx context.Context, username string) (*models.User, error)
	CreateUserFunc        func(ctx context.Context, user *models.User) error
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
		CreateUserFunc: func(ctx context.Context, user *models.User) error {
			return nil
		},
	}
}

func (mock *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return mock.GetUserByUsernameFunc(ctx, username)
}

func (mock *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return mock.CreateUserFunc(ctx, user)
}
