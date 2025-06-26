package authApp

import (
	"github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/domain/token"
	userDomain "github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/domain/user"
)

type AuthService interface {
	Register(email, username, password string) (tokens *token.TokenPayload, err error)
	Login(email, password string) (tokens *token.TokenPayload, err error)
	VerifyToken(accessToken string) (user *userDomain.User, err error)
	RefreshToken(refreshToken string) (tokens *token.TokenPayload, err error)
}
