package repository

import (
	"UserCrud/model"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type UniqueConstraintError struct {
	Field string
}

func (e *UniqueConstraintError) Error() string {
	return fmt.Sprintf("%s is already used by another user", e.Field)
}

func NewUniqueConstraintError(Field string) *UniqueConstraintError {
	return &UniqueConstraintError{Field: Field}
}

var ErrUserNotFound = errors.New("user not found")

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
		var err *pgconn.PgError
		if errors.As(result.Error, &err) && err.Code == "23505" {
			c := strings.Join(strings.Split(err.ConstraintName, "_")[2:], " ")
			return NewUniqueConstraintError(c)
		}
	}
	return result.Error
}

func (u *userRepository) UpdateUser(ctx context.Context, user model.User) error {
	result := u.db.WithContext(ctx).Model(&user).Updates(user)
	if result.Error != nil {
		var err *pgconn.PgError
		if errors.As(result.Error, &err) && err.Code == "23505" {
			c := strings.Join(strings.Split(err.ConstraintName, "_")[2:], " ")
			return NewUniqueConstraintError(c)
		}
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}

func (u *userRepository) DeleteUser(ctx context.Context, id uint) error {
	result := u.db.WithContext(ctx).Delete(&model.User{}, id)
	return result.Error
}

func (u *userRepository) GetUserById(ctx context.Context, id uint) (model.User, error) {
	var user model.User
	result := u.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.User{}, ErrUserNotFound
		}
	}
	return user, result.Error
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	result := u.db.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.User{}, ErrUserNotFound
		}
	}
	return user, result.Error
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}
