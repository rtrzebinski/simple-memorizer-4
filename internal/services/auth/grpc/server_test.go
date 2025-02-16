package grpc

import (
	"context"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"github.com/stretchr/testify/assert"
)

func TestServer_Register(t *testing.T) {
	ctx := context.Background()

	req := new(grpc.RegisterRequest)

	req.Email = "email"
	req.Name = "name"
	req.Password = "password"

	serviceMock := &ServiceMock{}
	serviceMock.On("Register", ctx, req.Name, req.Email, req.Password).Return("token", nil)

	server := NewServer(serviceMock)

	res, err := server.Register(ctx, req)
	assert.NoError(t, err)

	assert.Equal(t, "token", res.AccessToken)
}

func TestServer_Register_fail(t *testing.T) {
	ctx := context.Background()

	req := new(grpc.RegisterRequest)

	req.Email = "email"
	req.Name = "name"
	req.Password = "password"

	serviceMock := &ServiceMock{}
	serviceMock.On("Register", ctx, req.Name, req.Email, req.Password).Return("", assert.AnError)

	server := NewServer(serviceMock)

	_, err := server.Register(ctx, req)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &assert.AnError)
}

func TestServer_SignIn(t *testing.T) {
	ctx := context.Background()

	req := new(grpc.SignInRequest)

	req.Email = "email"
	req.Password = "password"

	serviceMock := &ServiceMock{}
	serviceMock.On("SignIn", ctx, req.Email, req.Password).Return("token", nil)

	server := NewServer(serviceMock)

	res, err := server.SignIn(ctx, req)
	assert.NoError(t, err)

	assert.Equal(t, "token", res.AccessToken)
}

func TestServer_SignIn_fail(t *testing.T) {
	ctx := context.Background()

	req := new(grpc.SignInRequest)

	req.Email = "email"
	req.Password = "password"

	serviceMock := &ServiceMock{}
	serviceMock.On("SignIn", ctx, req.Email, req.Password).Return("", assert.AnError)

	server := NewServer(serviceMock)

	_, err := server.SignIn(ctx, req)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &assert.AnError)
}
