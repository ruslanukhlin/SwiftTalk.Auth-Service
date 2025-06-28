package authGRPC

import (
	"context"

	authApp "github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/application/auth"
	"github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/domain/user"
	passwordDomain "github.com/ruslanukhlin/SwiftTalk.Auth-service/internal/domain/user/password"
	pb "github.com/ruslanukhlin/SwiftTalk.Common/gen/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrPasswordTooShort   = status.Error(codes.InvalidArgument, passwordDomain.ErrPasswordTooShort.Error())
	ErrPasswordEmpty      = status.Error(codes.InvalidArgument, passwordDomain.ErrPasswordEmpty.Error())
	ErrInvalidPassword    = status.Error(codes.InvalidArgument, passwordDomain.ErrInvalidPassword.Error())
	ErrEmailAlreadyExists = status.Error(codes.AlreadyExists, user.ErrEmailAlreadyExists.Error())
	ErrInvalidEmail       = status.Error(codes.InvalidArgument, user.ErrInvalidEmail.Error())
	ErrUserNameRequired   = status.Error(codes.InvalidArgument, user.ErrUserNameRequired.Error())
	ErrUserNameTooShort   = status.Error(codes.InvalidArgument, user.ErrUserNameTooShort.Error())
	ErrInternal           = status.Error(codes.Internal, "Внутренняя ошибка сервера")
	ErrUnauthorized       = status.Error(codes.Unauthenticated, "Не авторизован")
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
	tokens, err := h.authApp.Register(req.Email, req.Username, req.Password)
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
		case user.ErrUserNameRequired:
			return nil, ErrUserNameRequired
		case user.ErrUserNameTooShort:
			return nil, ErrUserNameTooShort
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
		case passwordDomain.ErrInvalidPassword, user.ErrUserNotFound:
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, ErrUnauthorized
	}

	user, err := h.authApp.VerifyToken(values[0])
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &pb.VerifyTokenResponse{
		IsValid:  true,
		UserUuid: user.UUID,
		Username: user.Username.Value,
		Email:    user.Email.Value,
	}, nil
}

func (h *UserGRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, ErrUnauthorized
	}

	refreshToken := values[0]
	tokens, err := h.authApp.RefreshToken(refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
