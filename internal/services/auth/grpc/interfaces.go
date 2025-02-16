package grpc

import "context"

type Service interface {
	Register(ctx context.Context, name, email, password string) (authToken string, err error)
	SignIn(ctx context.Context, email, password string) (authToken string, err error)
}
