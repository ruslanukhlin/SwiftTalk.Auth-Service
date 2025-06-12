package authApp

import (
	tokenDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/token"
	userDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user"
	passwordDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user/password"
)

var _ AuthService = &AuthApp{}

type AuthApp struct {
	userRepo userDomain.UserRepository
	passwordRepo passwordDomain.PasswordRepository
	tokenRepo tokenDomain.TokenRepository
}

func NewAuthApp(userRepo userDomain.UserRepository, passwordRepo passwordDomain.PasswordRepository, tokenRepo tokenDomain.TokenRepository) *AuthApp {
	return &AuthApp{
		userRepo: userRepo,
		passwordRepo: passwordRepo,
		tokenRepo: tokenRepo,
	}
}

func (a *AuthApp) Register(email, password string) (tokens *tokenDomain.TokenPayload, err error) {
	user, err := userDomain.NewUser(email, password, a.passwordRepo, a.userRepo)
	if err != nil {
		return nil, err
	}

	err = a.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	accessPayload := tokenDomain.NewAccessTokenClaim(user.UUID)
	refreshPayload := tokenDomain.NewRefreshTokenClaim(user.UUID)
	
	tokens, err = a.tokenRepo.CreateToken(accessPayload, refreshPayload)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (a *AuthApp) Login(email, password string) (tokens *tokenDomain.TokenPayload, err error) {
	user, err := a.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := passwordDomain.ComparePassword(password, user.Password.Value, a.passwordRepo); err != nil {
		return nil, err
	}
	
	accessPayload := tokenDomain.NewAccessTokenClaim(user.UUID)
	refreshPayload := tokenDomain.NewRefreshTokenClaim(user.UUID)
	
	tokens, err = a.tokenRepo.CreateToken(accessPayload, refreshPayload)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (a *AuthApp) VerifyToken(accessToken string) (user *userDomain.User, err error) {
	accessPayload, err := a.tokenRepo.ParseToken(accessToken)
	if err != nil {
		return nil, err
	}

	user, err = a.userRepo.GetUserByUUID(accessPayload.UUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthApp) RefreshToken(refreshToken string) (tokens *tokenDomain.TokenPayload, err error) {
	refreshParsed, err := a.tokenRepo.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}
	
	user, err := a.userRepo.GetUserByUUID(refreshParsed.UUID)
	if err != nil {
		return nil, err
	}

	accessPayload := tokenDomain.NewAccessTokenClaim(user.UUID)
	refreshPayload := tokenDomain.NewRefreshTokenClaim(user.UUID)
	
	tokens, err = a.tokenRepo.CreateToken(accessPayload, refreshPayload)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}