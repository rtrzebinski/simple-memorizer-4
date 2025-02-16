package auth

import "context"

type Reader interface {
	SignIn(ctx context.Context, email, password string) (name, userID string, err error)
}

type Writer interface {
	Register(ctx context.Context, name, email, password string) (userID string, err error)
}
