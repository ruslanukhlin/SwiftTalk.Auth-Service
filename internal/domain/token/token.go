package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
)

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type AccessTokenClaim struct {
	jwt.RegisteredClaims
	TokenType TokenType `json:"token_type"`
}

type RefreshTokenClaim struct {
	jwt.RegisteredClaims
	TokenType TokenType `json:"token_type"`
}

type TokenPayload struct {
	AccessToken  string
	RefreshToken string
}

func NewAccessTokenClaim(uuid string, cfg *config.Config) *AccessTokenClaim {
	now := time.Now()
	return &AccessTokenClaim{
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWT.Issuer,
			Subject:   uuid,
			Audience:  jwt.ClaimStrings{cfg.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.JWT.ExpiresAfter)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}

func NewRefreshTokenClaim(uuid string, cfg *config.Config) *RefreshTokenClaim {
	now := time.Now()
	return &RefreshTokenClaim{
		TokenType: RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWT.Issuer,
			Subject:   uuid,
			Audience:  jwt.ClaimStrings{cfg.JWT.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.JWT.RefreshExpiresAfter)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}