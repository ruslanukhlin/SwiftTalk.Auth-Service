package password

import "errors"

var (
	ErrPasswordEmpty    = errors.New("пароль не может быть пустым")
	ErrPasswordTooShort = errors.New("пароль должен быть не менее 8 символов")

	ErrInvalidPassword = errors.New("неверный email или пароль")
)

type Password struct {
	Value string
}

func NewPassword(value string) (*Password, error) {
	if value == "" {
		return nil, ErrPasswordEmpty
	}

	if len(value) < 8 {
		return nil, ErrPasswordTooShort
	}

	return &Password{
		Value: value,
	}, nil
}

func ComparePassword(password string, hash string, passwordRepo PasswordRepository) error {
	if err := passwordRepo.ComparePassword(password, hash); err != nil {
		return ErrInvalidPassword
	}

	return nil
}
