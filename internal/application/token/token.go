package application

import (
	domain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/token"
)

var _ domain.TokenService = &TokenApp{}

type TokenApp struct {
	tokenRepo domain.TokenRepository
}

func NewTokenApp(tokenRepo domain.TokenRepository) *TokenApp {
	return &TokenApp{
		tokenRepo: tokenRepo,
	}
}

func (a *TokenApp) CreateToken(uuid string) (*domain.TokenPayload, error) {
	accessPayload := domain.NewAccessTokenClaim(uuid)
	refreshPayload := domain.NewRefreshTokenClaim(uuid)

	return a.tokenRepo.CreateToken(accessPayload, refreshPayload)
}

func (a *TokenApp) ParseToken(token string) (*domain.AccessTokenClaim, error) {
	return a.tokenRepo.ParseToken(token)
}