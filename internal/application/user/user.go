package application

import (
	passwordDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/password"
	tokenDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/token"
	userDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user"
)

var _ userDomain.UserService = &UserApp{}

type UserApp struct {
	userRepo userDomain.UserRepository
	tokenRepo tokenDomain.TokenRepository
	passwordRepo passwordDomain.PasswordRepository
}

func NewUserApp(userRepo userDomain.UserRepository, tokenRepo tokenDomain.TokenRepository, passwordRepo passwordDomain.PasswordRepository) *UserApp {
	return &UserApp{
		userRepo: userRepo,
		tokenRepo: tokenRepo,
		passwordRepo: passwordRepo,
	}
}

func (a *UserApp) Register(user *userDomain.User) error {
	hashedPassword, err := a.passwordRepo.HashPassword(user.Password)

	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return a.userRepo.CreateUser(user)
}

func (a *UserApp) Login(email, password string) (*userDomain.User, error) {
	user, err := a.userRepo.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	if !a.passwordRepo.ComparePassword(password, user.Password) {
		return nil, err
	}

	return user, nil
}

func (a *UserApp) VerifyToken(token string) (*userDomain.User, error) {
	claims, err := a.tokenRepo.ParseToken(token)
	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.GetUserByUUID(claims.UUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UserApp) RefreshToken(refreshToken string) (*userDomain.User, error) {
	claims, err := a.tokenRepo.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.GetUserByUUID(claims.UUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}