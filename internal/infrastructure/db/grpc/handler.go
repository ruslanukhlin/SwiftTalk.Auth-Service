package userGRPC

import (
	"context"

	application "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/user"
	domain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user"
	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	userApp *application.UserApp
}

func NewUserGRPCHandler(userApp *application.UserApp) *UserGRPCHandler {
	return &UserGRPCHandler{
		userApp: userApp,
	}
}

func (h *UserGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := domain.NewUser(req.Email, req.Password)

	if err := h.userApp.Register(user); err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при регистрации пользователя: %v", err)
	}

	return &pb.RegisterResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}, nil
}