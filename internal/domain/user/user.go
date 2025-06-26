package user

import (
	"time"

	"github.com/google/uuid"
	passwordDomain "github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/domain/user/password"
)

type User struct {
	UUID      string
	Email     Email
	Username  UserName
	Password  passwordDomain.HashPassword
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, username, password string, passwordRepo passwordDomain.PasswordRepository, userRepo UserRepository) (*User, error) {
	emailValid, err := NewEmail(email, userRepo)
	if err != nil {
		return nil, err
	}

	passwordValid, err := passwordDomain.NewPassword(password)
	if err != nil {
		return nil, err
	}

	hashPassword, err := passwordDomain.NewHashPassword(*passwordValid, passwordRepo)
	if err != nil {
		return nil, err
	}

	usernameValid, err := NewUserName(username)
	if err != nil {
		return nil, err
	}

	user := &User{
		UUID:      uuid.New().String(),
		Email:     *emailValid,
		Username:  *usernameValid,
		Password:  *hashPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}
