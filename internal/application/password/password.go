package application

import (
	passwordDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/password"
)

var _ passwordDomain.PasswordService = &PasswordApp{}

type PasswordApp struct {
	passwordRepo passwordDomain.PasswordRepository
}

func NewPasswordApp(passwordRepo passwordDomain.PasswordRepository) *PasswordApp {
	return &PasswordApp{
		passwordRepo: passwordRepo,
	}
}

func (a *PasswordApp) HashPassword(password string) (string, error) {
	return a.passwordRepo.HashPassword(password)
}

func (a *PasswordApp) ComparePassword(password, hash string) bool {
	return a.passwordRepo.ComparePassword(password, hash)
}