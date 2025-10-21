package grpc

import (
	"context"

	gengrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type AuthServiceClientMock struct {
	mock.Mock
}

func NewAuthServiceClientMock() *AuthServiceClientMock {
	return &AuthServiceClientMock{}
}

func (m *AuthServiceClientMock) Register(ctx context.Context, in *gengrpc.RegisterRequest, opts ...grpc.CallOption) (*gengrpc.RegisterResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gengrpc.RegisterResponse), args.Error(1)
}

func (m *AuthServiceClientMock) SignIn(ctx context.Context, in *gengrpc.SignInRequest, opts ...grpc.CallOption) (*gengrpc.SignInResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gengrpc.SignInResponse), args.Error(1)
}
