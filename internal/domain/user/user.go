package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, password string) *User {
	return &User{
		UUID:      uuid.New().String(),
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}