package service

import (
	"UserCrud/model"
	"UserCrud/repository"
	"UserCrud/util"
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"reflect"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, user model.User) error
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id uint) error
	GetById(ctx context.Context, id uint) (model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	jwt      util.JwtUtil
	logger   *zap.Logger
}

func (u *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong password")
	}
	token, err := u.jwt.CreateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userService) Register(ctx context.Context, user model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("failed to generate password", zap.Error(err))
		return err
	}
	user.Password = string(hash)
	err = u.userRepo.CreateUser(ctx, user)
	return err
}

func (u *userService) Update(ctx context.Context, user model.User) error {
	if !reflect.ValueOf(user.Password).IsZero() {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			u.logger.Error("failed to generate password", zap.Error(err))
			return err
		}
		user.Password = string(hash)
	}
	err := u.userRepo.UpdateUser(ctx, user)
	return err
}

func (u *userService) Delete(ctx context.Context, id uint) error {
	return u.userRepo.DeleteUser(ctx, id)
}

func (u *userService) GetById(ctx context.Context, id uint) (model.User, error) {
	return u.userRepo.GetUserById(ctx, id)
}

func NewUserService(userRepo repository.UserRepository, jwt util.JwtUtil, logger *zap.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		jwt:      jwt,
		logger:   logger,
	}
}
