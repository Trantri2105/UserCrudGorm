package repository

import (
	"UserCrud/model"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) error
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, id uint) error
	GetUserById(ctx context.Context, id uint) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func (u *userRepository) CreateUser(ctx context.Context, user model.User) error {
	result := u.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		u.logger.Error("Failed to create user", zap.Error(result.Error))
	}
	return result.Error
}

func (u *userRepository) UpdateUser(ctx context.Context, user model.User) error {
	result := u.db.WithContext(ctx).Model(&user).Updates(user)
	if result.Error != nil {
		u.logger.Error("Failed to update user", zap.Error(result.Error))
	}
	return result.Error
}

func (u *userRepository) DeleteUser(ctx context.Context, id uint) error {
	result := u.db.WithContext(ctx).Delete(&model.User{}, id)
	if result.Error != nil {
		u.logger.Error("Failed to delete user", zap.Error(result.Error))
	}
	return result.Error
}

func (u *userRepository) GetUserById(ctx context.Context, id uint) (model.User, error) {
	var user model.User
	result := u.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		u.logger.Error("Failed to get user", zap.Error(result.Error))
	}
	return user, result.Error
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	result := u.db.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		u.logger.Error("Failed to get user", zap.Error(result.Error))
	}
	return user, result.Error
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}
