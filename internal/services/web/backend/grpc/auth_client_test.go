package grpc

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"

	gengrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"github.com/stretchr/testify/assert"
)

func TestAuthAPIClient_Register(t *testing.T) {
	ctx := context.Background()

	clientMock := NewAuthServiceClientMock()
	clientMock.On("Register", ctx, &gengrpc.RegisterRequest{
		Name:     "name",
		Email:    "email",
		Password: "password",
	}, mock.Anything).Return(&gengrpc.RegisterResponse{
		AccessToken: "accessToken",
	}, nil)

	authClient := NewAuthClient(clientMock)

	accessToken, err := authClient.Register(ctx, "name", "email", "password")

	assert.NoError(t, err)
	assert.Equal(t, "accessToken", accessToken)
}

func TestAuthAPIClient_SignIn(t *testing.T) {
	ctx := context.Background()

	clientMock := NewAuthServiceClientMock()
	clientMock.On("SignIn", ctx, &gengrpc.SignInRequest{
		Email:    "email",
		Password: "password",
	}, mock.Anything).Return(&gengrpc.SignInResponse{
		AccessToken: "accessToken",
	}, nil)

	authClient := NewAuthClient(clientMock)

	accessToken, err := authClient.SignIn(ctx, "email", "password")

	assert.NoError(t, err)
	assert.Equal(t, "accessToken", accessToken)
}
