package jwtRepo

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/golang-jwt/jwt/v5"
	tokenDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/token"
)

var _ tokenDomain.TokenRepository = &JWTTokenRepository{}

type JWTTokenRepository struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTTokenRepository(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *JWTTokenRepository {
	return &JWTTokenRepository{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (r *JWTTokenRepository) CreateToken(accessPayload *tokenDomain.AccessTokenClaim, refreshPayload *tokenDomain.RefreshTokenClaim) (*tokenDomain.TokenPayload, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessPayload)
	signedAccessToken, err := accessToken.SignedString(r.privateKey)

	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshPayload)
	signedRefreshToken, err := refreshToken.SignedString(r.privateKey)	
	if err != nil {
		return nil, err
	}

	return &tokenDomain.TokenPayload{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (r *JWTTokenRepository) ParseToken(token string) (*tokenDomain.AccessTokenClaim, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &tokenDomain.AccessTokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
		}
		return r.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*tokenDomain.AccessTokenClaim); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (r *JWTTokenRepository) GetJWKS() ([]byte, error) {
	// Convert modulus to Base64URL
	nBytes := r.publicKey.N.Bytes()
	nBase64 := base64.RawURLEncoding.EncodeToString(nBytes)

	// Convert exponent to Base64URL
	eBytes := big.NewInt(int64(r.publicKey.E)).Bytes()
	eBase64 := base64.RawURLEncoding.EncodeToString(eBytes)

	jwks := struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			Use string `json:"use"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}{
		Keys: []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			Use string `json:"use"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		}{
			{
				Kty: "RSA",
				Kid: "1",
				Use: "sig",
				Alg: "RS256",
				N:   nBase64,
				E:   eBase64,
			},
		},
	}

	return json.Marshal(jwks)
}