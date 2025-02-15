package grpc

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
)

type Server struct {
	grpc.UnimplementedAuthServiceServer
}

func (s *Server) SignUp(_ context.Context, _ *grpc.SignUpRequest) (*grpc.SignUpResponse, error) {
	return &grpc.SignUpResponse{
		AccessToken: "accessToken",
	}, nil
}

func (s *Server) SignIn(_ context.Context, _ *grpc.SignInRequest) (*grpc.SignInResponse, error) {
	return &grpc.SignInResponse{
		AccessToken: "accessToken",
	}, nil
}
