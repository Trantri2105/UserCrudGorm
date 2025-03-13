package util

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"os"
	"time"
)

type JwtUtil interface {
	CreateToken(userId uint) (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

const jwtTokenExpTime = 60 * time.Minute

var ErrInvalidToken = errors.New("invalid token")

type jwtUtil struct {
	logger *zap.Logger
}

func (j *jwtUtil) CreateToken(userId uint) (string, error) {
	expireTime := time.Now().Add(jwtTokenExpTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    expireTime,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		j.logger.Error("Failed to create token", zap.Error(err))
		return "", err
	}
	return tokenString, nil
}

func (j *jwtUtil) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		j.logger.Error("Failed to parse token", zap.Error(err))
		return nil, ErrInvalidToken
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}

func NewJwtUtils(logger *zap.Logger) JwtUtil {
	return &jwtUtil{
		logger: logger,
	}
}
