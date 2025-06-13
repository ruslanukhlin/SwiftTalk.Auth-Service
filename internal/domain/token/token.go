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
			Issuer:    "auth-service",
			Subject:   uuid,
			Audience:  jwt.ClaimStrings{"swift-talk"},
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
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
			Issuer:    "auth-service",
			Subject:   uuid,
			Audience:  jwt.ClaimStrings{"swift-talk"},
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}