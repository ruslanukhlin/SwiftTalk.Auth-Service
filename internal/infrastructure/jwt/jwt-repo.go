package jwtRepo

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/token"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
)

var _ token.TokenRepository = &JWTTokenRepository{}

type JWTTokenRepository struct {
	cfg *config.JWTConfig
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	token.AccessTokenClaim
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	token.RefreshTokenClaim
}

func NewJWTTokenRepository(cfg *config.JWTConfig) *JWTTokenRepository {
	return &JWTTokenRepository{cfg: cfg}
}

func (r *JWTTokenRepository) CreateToken(accessPayload *token.AccessTokenClaim, refreshPayload *token.RefreshTokenClaim) (*token.TokenPayload, error) {
	expiresAt := time.Now().Add(r.cfg.ExpiresAfter)
	refreshExpiresAt := time.Now().Add(r.cfg.RefreshExpiresAfter)
	
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		AccessTokenClaim: *accessPayload,
	}).SignedString([]byte(r.cfg.SecretKey))

	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
		},
		RefreshTokenClaim: *refreshPayload,
	}).SignedString([]byte(r.cfg.SecretKey))

	if err != nil {
		return nil, err
	}

	return &token.TokenPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (r *JWTTokenRepository) ParseToken(token string) (*token.AccessTokenClaim, error) {
	claims := &AccessTokenClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.cfg.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return &claims.AccessTokenClaim, nil
}