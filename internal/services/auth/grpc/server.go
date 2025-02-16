package grpc

import (
	"context"
	"log/slog"

	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	grpc.UnimplementedAuthServiceServer
	s Service
}

func NewServer(s Service) *Server {
	return &Server{s: s}
}

func (s *Server) Register(ctx context.Context, req *grpc.RegisterRequest) (*grpc.RegisterResponse, error) {
	accessToken, err := s.s.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		slog.Warn("failed to register user", "error", err)

		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &grpc.RegisterResponse{AccessToken: accessToken}, nil
}

func (s *Server) SignIn(ctx context.Context, req *grpc.SignInRequest) (*grpc.SignInResponse, error) {
	accessToken, err := s.s.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		slog.Warn("failed to sign in user", "error", err)

		return nil, status.Errorf(codes.Internal, "failed to sign in user: %v", err)
	}

	return &grpc.SignInResponse{AccessToken: accessToken}, nil
}
