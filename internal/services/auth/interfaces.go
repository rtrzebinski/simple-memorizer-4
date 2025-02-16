package auth

import "context"

type Reader interface {
	FetchUser(ctx context.Context, email string) (name, password, userID string, err error)
}

type Writer interface {
	StoreUser(ctx context.Context, name, email, password string) (userID string, err error)
}
