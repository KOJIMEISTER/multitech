package storage

import (
	"context"
	"errors"
	"multitech/internal/models"
	"strings"

	"gorm.io/gorm"
)

type gormUserRepository struct {
	*gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db}
}

func (userRepo *gormUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := userRepo.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (userRepo *gormUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	err := userRepo.WithContext(ctx).Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ErrUserExists
		}
		return err
	}
	return nil
}
