package passwordRepo

import (
	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user/password"
	"golang.org/x/crypto/bcrypt"
)

var _ password.PasswordRepository = &passwordRepo{}
	
type passwordRepo struct {}

func NewPasswordRepo() *passwordRepo {
	return &passwordRepo{}
}

func (r *passwordRepo) HashPassword(password password.Password) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.Value), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}	
	
	return string(hashedPassword), nil	
}

func (r *passwordRepo) ComparePassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}