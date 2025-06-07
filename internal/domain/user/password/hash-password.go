package password

import "errors"

var (
	ErrHashPassword = errors.New("ошибка хеширования пароля")
)

type HashPassword struct {
	Value string
}

func NewHashPassword(password Password, passwordRepo PasswordRepository) (*HashPassword, error) {
	hashPassword, err := passwordRepo.HashPassword(password)

	if err != nil {
		return nil, ErrHashPassword
	}

	return &HashPassword{
		Value: hashPassword,
	}, nil
}