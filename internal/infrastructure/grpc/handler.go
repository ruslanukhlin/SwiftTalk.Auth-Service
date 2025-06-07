package userGRPC

import (
	"context"

	userApp "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/token"
	tokenApp "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/user"
	userDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user"
	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRegisterUser = status.Error(codes.Internal, "ошибка при регистрации пользователя")
	ErrCreateToken = status.Error(codes.Internal, "ошибка при создании токена")
	ErrHashPassword = status.Error(codes.Internal, "ошибка при хешировании пароля")
	ErrLogin = status.Error(codes.Unauthenticated, "неверный email или пароль")
	ErrVerifyToken = status.Error(codes.Unauthenticated, "неверный токен")
	ErrRefreshToken = status.Error(codes.Unauthenticated, "неверный refresh токен")
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
	user := userDomain.NewUser(req.Email, req.Password)
	if err := h.userApp.Register(user); err != nil {
		return nil, ErrRegisterUser
	}

	tokens, err := h.tokenApp.CreateToken(user.UUID)
	if err != nil {
		return nil, ErrCreateToken
	}

	return &pb.RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *UserGRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := h.userApp.Login(req.Email, req.Password)
	if err != nil {
		return nil, ErrLogin
	}

	tokens, err := h.tokenApp.CreateToken(user.UUID)
	if err != nil {
		return nil, ErrCreateToken
	}

	return &pb.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *UserGRPCHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	user, err := h.userApp.VerifyToken(req.AccessToken)

	if err != nil {
		return nil, ErrVerifyToken
	}

	return &pb.VerifyTokenResponse{
		IsValid: true,
		UserId: user.UUID,
	}, nil
}

func (h *UserGRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	user, err := h.userApp.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, ErrRefreshToken
	}

	tokens, err := h.tokenApp.CreateToken(user.UUID)
	if err != nil {
		return nil, ErrCreateToken
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}