package token

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type AccessTokenClaim struct {
	UUID string
	TokenType TokenType
}

type RefreshTokenClaim struct {
	UUID string
	TokenType TokenType
}

type TokenPayload struct {
	AccessToken  string
	RefreshToken string
}

func NewAccessTokenClaim(uuid string) *AccessTokenClaim {
	return &AccessTokenClaim{
		UUID:      uuid,
		TokenType: AccessToken,
	}
}	

func NewRefreshTokenClaim(uuid string) *RefreshTokenClaim {
	return &RefreshTokenClaim{
		UUID:      uuid,
		TokenType: RefreshToken,
	}
}