package user

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidEmail = errors.New("не корректный email")
	ErrEmailAlreadyExists = errors.New("email уже существует")
)

type Email struct {
	Value string `validate:"required,email"`
}

func NewEmail(value string, userRepo UserRepository) (*Email, error) {
	validate := validator.New()

	isExists, err := userRepo.IsEmailExists(value)
	if err != nil {
		return nil, err
	}
	if isExists {
		return nil, ErrEmailAlreadyExists
	}

	if err := validate.Struct(Email{Value: value}); err != nil {
		return nil, ErrInvalidEmail
	}

	return &Email{Value: value}, nil
}