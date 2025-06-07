package userGRPC

import (
	"context"
	"log"

	userApp "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/token"
	tokenApp "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/user"
	domain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user"
	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRegisterUser = status.Error(codes.Internal, "ошибка при регистрации пользователя")
	ErrCreateToken = status.Error(codes.Internal, "ошибка при создании токена")
)

type UserGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	userApp *tokenApp.UserApp
	tokenApp *userApp.TokenApp
}

func NewUserGRPCHandler(userApp *tokenApp.UserApp, tokenApp *userApp.TokenApp) *UserGRPCHandler {
	return &UserGRPCHandler{
		userApp: userApp,
		tokenApp: tokenApp,
	}
}

func (h *UserGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := domain.NewUser(req.Email, req.Password)

	if err := h.userApp.Register(user); err != nil {
		log.Println("Error registering user:", err)
		return nil, ErrRegisterUser
	}

	tokens, err := h.tokenApp.CreateToken(user.UUID)
	if err != nil {
		log.Println("Error creating token:", err)
		return nil, ErrCreateToken
	}

	return &pb.RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}