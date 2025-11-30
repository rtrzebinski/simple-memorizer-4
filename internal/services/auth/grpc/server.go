package grpc

import (
	"context"
	"log/slog"

	gengrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	gengrpc.UnimplementedAuthServiceServer
	s *auth.Service
}

func NewServer(s *auth.Service) *Server {
	return &Server{s: s}
}

func (s *Server) Register(ctx context.Context, req *gengrpc.RegisterRequest) (*gengrpc.Tokens, error) {
	t, err := s.s.Register(ctx, req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		slog.Warn("failed to register user", "error", err)

		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	slog.Info("Auth service - user registered")

	return &gengrpc.Tokens{
		AccessToken:      t.AccessToken,
		IdToken:          t.IDToken,
		ExpiresIn:        int32(t.ExpiresIn),
		RefreshExpiresIn: int32(t.RefreshExpiresIn),
		RefreshToken:     t.RefreshToken,
		TokenType:        t.TokenType,
	}, nil
}

func (s *Server) SignIn(ctx context.Context, req *gengrpc.SignInRequest) (*gengrpc.Tokens, error) {
	t, err := s.s.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		slog.Warn("failed to sign in user", "error", err)

		return nil, status.Errorf(codes.Internal, "failed to sign in user: %v", err)
	}

	slog.Info("Auth service - user signed in")

	return &gengrpc.Tokens{
		AccessToken:      t.AccessToken,
		IdToken:          t.IDToken,
		ExpiresIn:        int32(t.ExpiresIn),
		RefreshExpiresIn: int32(t.RefreshExpiresIn),
		RefreshToken:     t.RefreshToken,
		TokenType:        t.TokenType,
	}, nil
}

func (s *Server) Refresh(ctx context.Context, req *gengrpc.RefreshRequest) (*gengrpc.Tokens, error) {
	t, err := s.s.Refresh(ctx, req.RefreshToken)
	if err != nil {
		slog.Warn("failed to refresh token", "error", err)

		return nil, status.Errorf(codes.Internal, "failed to refresh token: %v", err)
	}

	slog.Info("Auth service - token refreshed")

	return &gengrpc.Tokens{
		AccessToken:      t.AccessToken,
		IdToken:          t.IDToken,
		ExpiresIn:        int32(t.ExpiresIn),
		RefreshExpiresIn: int32(t.RefreshExpiresIn),
		RefreshToken:     t.RefreshToken,
		TokenType:        t.TokenType,
	}, nil
}

func (s *Server) Revoke(ctx context.Context, req *gengrpc.RevokeRequest) (*emptypb.Empty, error) {
	err := s.s.Revoke(ctx, req.RefreshToken)
	if err != nil {
		slog.Warn("failed to revoke token", "error", err)

		return nil, status.Errorf(codes.Internal, "failed to revoke token: %v", err)
	}

	slog.Info("Auth service - token revoked")

	return &emptypb.Empty{}, nil
}
