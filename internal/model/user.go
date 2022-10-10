package model

import (
	"time"
)

type User struct {
	ID           int
	Password     string
	PasswordHash string
	Username     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
