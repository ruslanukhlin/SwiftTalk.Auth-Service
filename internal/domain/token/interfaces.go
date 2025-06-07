package domain

type TokenService interface {
	CreateToken(uuid string) (*TokenPayload, error)
	ParseToken(token string) (*AccessTokenClaim, error)
}

type TokenRepository interface {
	CreateToken(accessPayload *AccessTokenClaim, refreshPayload *RefreshTokenClaim) (*TokenPayload, error)
	ParseToken(token string) (*AccessTokenClaim, error)
}