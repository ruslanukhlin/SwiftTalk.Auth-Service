package jwtRepo

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	domain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/token"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/pkg/config"
)

var _ domain.TokenRepository = &JWTTokenRepository{}

type JWTTokenRepository struct {
	cfg *config.JWTConfig
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	domain.AccessTokenClaim
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	domain.RefreshTokenClaim
}

func NewJWTTokenRepository(cfg *config.JWTConfig) *JWTTokenRepository {
	return &JWTTokenRepository{cfg: cfg}
}

func (r *JWTTokenRepository) CreateToken(accessPayload *domain.AccessTokenClaim, refreshPayload *domain.RefreshTokenClaim) (*domain.TokenPayload, error) {
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

	return &domain.TokenPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}