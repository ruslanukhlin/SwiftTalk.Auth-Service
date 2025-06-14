package authGRPC

import (
	"context"

	authApp "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/application/auth"
	"github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user"
	passwordDomain "github.com/ruslanukhlin/SwiftTalk.auth-service/internal/domain/user/password"
	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrPasswordTooShort   = status.Error(codes.InvalidArgument, passwordDomain.ErrPasswordTooShort.Error())
	ErrPasswordEmpty      = status.Error(codes.InvalidArgument, passwordDomain.ErrPasswordEmpty.Error())
	ErrInvalidPassword    = status.Error(codes.InvalidArgument, passwordDomain.ErrInvalidPassword.Error())
	ErrEmailAlreadyExists = status.Error(codes.AlreadyExists, user.ErrEmailAlreadyExists.Error())
	ErrInvalidEmail       = status.Error(codes.InvalidArgument, user.ErrInvalidEmail.Error())
	ErrInternal           = status.Error(codes.Internal, "Внутренняя ошибка сервера")
)

type UserGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	authApp authApp.AuthService
}

func NewUserGRPCHandler(authApp authApp.AuthService) *UserGRPCHandler {
	return &UserGRPCHandler{
		authApp: authApp,
	}
}

func (h *UserGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	tokens, err := h.authApp.Register(req.Email, req.Password)
	if err != nil {
		switch err {
		case user.ErrInvalidEmail:
			return nil, ErrInvalidEmail
		case user.ErrEmailAlreadyExists:
			return nil, ErrEmailAlreadyExists
		case passwordDomain.ErrPasswordTooShort:
			return nil, ErrPasswordTooShort
		case passwordDomain.ErrPasswordEmpty:
			return nil, ErrPasswordEmpty
		default:
			return nil, ErrInternal
		}
	}

	return &pb.RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *UserGRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	tokens, err := h.authApp.Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case passwordDomain.ErrInvalidPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, ErrInternal
		}
	}

	return &pb.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *UserGRPCHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	user, err := h.authApp.VerifyToken(req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &pb.VerifyTokenResponse{
		IsValid:  true,
		UserUuid: user.UUID,
	}, nil
}

func (h *UserGRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	tokens, err := h.authApp.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
