package auth

import (
	"context"

	gengrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcClient interface {
	Register(ctx context.Context, in *gengrpc.RegisterRequest, opts ...grpc.CallOption) (*gengrpc.Tokens, error)
	SignIn(ctx context.Context, in *gengrpc.SignInRequest, opts ...grpc.CallOption) (*gengrpc.Tokens, error)
	Refresh(ctx context.Context, in *gengrpc.RefreshRequest, opts ...grpc.CallOption) (*gengrpc.Tokens, error)
	Revoke(ctx context.Context, in *gengrpc.RevokeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}
