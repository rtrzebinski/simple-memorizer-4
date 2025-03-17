package grpc

import (
	"context"
	"fmt"

	gengrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
)

type AuthClient struct {
	c AuthServiceClient
}

func NewAuthClient(c AuthServiceClient) *AuthClient {
	return &AuthClient{c: c}
}

func (a *AuthClient) Register(ctx context.Context, name, email, password string) (accessToken string, err error) {
	req := gengrpc.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}

	res, err := a.c.Register(ctx, &req)

	if err != nil {
		return "", fmt.Errorf("register user: %w", err)
	}

	return res.AccessToken, nil
}

func (a *AuthClient) SignIn(ctx context.Context, email, password string) (accessToken string, err error) {
	req := gengrpc.SignInRequest{
		Email:    email,
		Password: password,
	}

	res, err := a.c.SignIn(ctx, &req)

	if err != nil {
		return "", fmt.Errorf("sign in user: %w", err)
	}

	return res.AccessToken, nil
}
