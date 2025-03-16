package grpc

import (
	"context"

	gengrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"google.golang.org/grpc"
)

type AuthServiceClient interface {
	Register(ctx context.Context, in *gengrpc.RegisterRequest, opts ...grpc.CallOption) (*gengrpc.RegisterResponse, error)
	SignIn(ctx context.Context, in *gengrpc.SignInRequest, opts ...grpc.CallOption) (*gengrpc.SignInResponse, error)
}
