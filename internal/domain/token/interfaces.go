package token

type TokenRepository interface {
	CreateToken(accessPayload *AccessTokenClaim, refreshPayload *RefreshTokenClaim) (*TokenPayload, error)
	ParseToken(token string) (*AccessTokenClaim, error)
	GetJWKS() ([]byte, error)
}