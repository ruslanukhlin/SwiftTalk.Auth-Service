package cryptRepo

import (
	passwordDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/password"
	"golang.org/x/crypto/bcrypt"
)

var _ passwordDomain.PasswordRepository = &CryptRepository{}

type CryptRepository struct {}

func NewCryptRepository() *CryptRepository {
	return &CryptRepository{}
}

func (r *CryptRepository) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (r *CryptRepository) ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	
	return err == nil
}