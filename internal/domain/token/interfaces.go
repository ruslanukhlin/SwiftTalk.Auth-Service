package domain

type TokenService interface {
	CreateToken(uuid string) (*TokenPayload, error)
}

type TokenRepository interface {
	CreateToken(accessPayload *AccessTokenClaim, refreshPayload *RefreshTokenClaim) (*TokenPayload, error)
}