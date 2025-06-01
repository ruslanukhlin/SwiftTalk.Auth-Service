package application

import (
	domain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user"
)

var _ domain.UserService = &UserApp{}

type UserApp struct {
	domain.UserRepository
}

func NewUserApp(userRepo domain.UserRepository) *UserApp {
	return &UserApp{
		UserRepository: userRepo,
	}
}

func (a *UserApp) Register(user *domain.User) error {
	return a.UserRepository.CreateUser(user)
}