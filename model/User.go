package model

import (
	"time"
)

type User struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	FirstName   string
	LastName    string
	Email       string `gorm:"unique"`
	Password    string
	PhoneNumber string `gorm:"unique"`
	Gender      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
