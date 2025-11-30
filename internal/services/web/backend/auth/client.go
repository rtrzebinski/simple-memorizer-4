package auth

import (
	"context"
	"fmt"

	gengrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type Client struct {
	grpcClient GrpcClient
}

func NewClient(grpcClient GrpcClient) *Client {
	return &Client{grpcClient: grpcClient}
}

func (c *Client) Register(ctx context.Context, firstName, lastName, email, password string) (backend.Tokens, error) {
	req := gengrpc.RegisterRequest{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}

	res, err := c.grpcClient.Register(ctx, &req)

	if err != nil {
		return backend.Tokens{}, fmt.Errorf("register user: %w", err)
	}

	return backend.Tokens{
		AccessToken:      res.AccessToken,
		IDToken:          res.IdToken,
		ExpiresIn:        int(res.ExpiresIn),
		RefreshExpiresIn: int(res.RefreshExpiresIn),
		RefreshToken:     res.RefreshToken,
		TokenType:        res.TokenType,
	}, nil
}

func (c *Client) SignIn(ctx context.Context, email, password string) (backend.Tokens, error) {
	req := gengrpc.SignInRequest{
		Email:    email,
		Password: password,
	}

	res, err := c.grpcClient.SignIn(ctx, &req)

	if err != nil {
		return backend.Tokens{}, fmt.Errorf("sign in user: %w", err)
	}

	return backend.Tokens{
		AccessToken:      res.AccessToken,
		IDToken:          res.IdToken,
		ExpiresIn:        int(res.ExpiresIn),
		RefreshExpiresIn: int(res.RefreshExpiresIn),
		RefreshToken:     res.RefreshToken,
		TokenType:        res.TokenType,
	}, nil
}

func (c *Client) Refresh(ctx context.Context, refreshToken string) (backend.Tokens, error) {
	req := gengrpc.RefreshRequest{
		RefreshToken: refreshToken,
	}

	res, err := c.grpcClient.Refresh(ctx, &req)

	if err != nil {
		return backend.Tokens{}, fmt.Errorf("refresh token: %w", err)
	}

	return backend.Tokens{
		AccessToken:      res.AccessToken,
		IDToken:          res.IdToken,
		ExpiresIn:        int(res.ExpiresIn),
		RefreshExpiresIn: int(res.RefreshExpiresIn),
		RefreshToken:     res.RefreshToken,
		TokenType:        res.TokenType,
	}, nil
}

func (c *Client) Revoke(ctx context.Context, refreshToken string) error {
	req := gengrpc.RevokeRequest{
		RefreshToken: refreshToken,
	}

	_, err := c.grpcClient.Revoke(ctx, &req)

	if err != nil {
		return fmt.Errorf("refresh token: %w", err)
	}

	return nil
}
