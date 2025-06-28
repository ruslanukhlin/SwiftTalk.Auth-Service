package bff

import (
	"context"

	"google.golang.org/grpc/metadata"

	pb "github.com/ruslanukhlin/SwiftTalk.Common/gen/auth"
)

type AuthService struct {
	client pb.AuthServiceClient
}

func NewAuthService(client pb.AuthServiceClient) *AuthService {
	return &AuthService{
		client: client,
	}
}

func (s *AuthService) Register(ctx context.Context, payload *RegisterPayload) (*TokenResponse, error) {
	response, err := s.client.Register(ctx, &pb.RegisterRequest{
		Email:    payload.Email,
		Username: payload.Username,
		Password: payload.Password,
	})
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, payload *LoginPayload) (*TokenResponse, error) {
	response, err := s.client.Login(ctx, &pb.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, payload *RefreshTokenPayload) (*TokenResponse, error) {
	md := metadata.New(map[string]string{"authorization": payload.RefreshToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := s.client.RefreshToken(ctx, &pb.RefreshTokenRequest{})
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, payload *VerifyTokenPayload) (*VerifyTokenResponse, error) {
	md := metadata.New(map[string]string{"authorization": payload.AccessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := s.client.VerifyToken(ctx, &pb.VerifyTokenRequest{})
	if err != nil {
		return nil, err
	}

	return &VerifyTokenResponse{
		IsValid:  response.IsValid,
		UserUUID: response.UserUuid,
		Username: response.Username,
		Email:    response.Email,
	}, nil
}
