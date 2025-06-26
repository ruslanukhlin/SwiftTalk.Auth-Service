package user

import "errors"

type UserName struct {
	Value string
}

var (
	ErrUserNameRequired = errors.New("никнейм не может быть пустым")
	ErrUserNameTooShort = errors.New("никнейм должен быть не менее 3 символов")
)

func NewUserName(value string) (*UserName, error) {
	if value == "" {
		return nil, ErrUserNameRequired
	}

	if len(value) < 4 {
		return nil, ErrUserNameTooShort
	}

	return &UserName{Value: value}, nil
}
