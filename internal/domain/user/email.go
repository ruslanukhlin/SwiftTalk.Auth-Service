package user

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidEmail = errors.New("не корректный email")
)

type Email struct {
	Value string `validate:"required,email"`
}

func NewEmail(value string) (*Email, error) {
	validate := validator.New()

	if err := validate.Struct(Email{Value: value}); err != nil {
		return nil, ErrInvalidEmail
	}

	return &Email{Value: value}, nil
}