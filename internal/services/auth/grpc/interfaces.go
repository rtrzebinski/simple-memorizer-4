package grpc

import "context"

type Service interface {
	Register(ctx context.Context, name, email, password string) (accessToken string, err error)
	SignIn(ctx context.Context, email, password string) (accessToken string, err error)
}
