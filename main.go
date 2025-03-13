package main

import (
	"UserCrud/handler"
	"UserCrud/middleware"
	"UserCrud/model"
	"UserCrud/repository"
	"UserCrud/service"
	"UserCrud/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func dbConnect(logger *zap.Logger) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}
	return db
}

func loadEnvVariable(logger *zap.Logger) {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}
}

func NewLogger() *zap.Logger {
	encodeConfig := zap.NewProductionConfig()
	encodeConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	syncConsole := zapcore.AddSync(os.Stderr)

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encodeConfig.EncoderConfig), zapcore.NewMultiWriteSyncer(syncConsole), zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())

	return logger
}

func main() {
	logger := NewLogger()
	loadEnvVariable(logger)
	db := dbConnect(logger)
	userRepo := repository.NewUserRepository(db, logger)
	jwt := util.NewJwtUtils(logger)
	userService := service.NewUserService(userRepo, jwt, logger)
	authMiddleware := middleware.NewAuthMiddleware(jwt)

	r := gin.Default()
	handler.AddUserHandler(userService, authMiddleware, r)
	err := r.Run()
	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
