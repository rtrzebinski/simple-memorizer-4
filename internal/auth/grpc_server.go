package auth

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
)

type GrpcServer struct {
	grpc.UnimplementedAuthServiceServer
}

func (s *GrpcServer) SignUp(_ context.Context, _ *grpc.SignUpRequest) (*grpc.SignUpResponse, error) {
	return &grpc.SignUpResponse{
		UserId:      "userID",
		AccessToken: "accessToken",
	}, nil
}

func (s *GrpcServer) SignIn(_ context.Context, _ *grpc.SignInRequest) (*grpc.SignInResponse, error) {
	return &grpc.SignInResponse{
		UserId:      "userID",
		AccessToken: "accessToken",
	}, nil
}
