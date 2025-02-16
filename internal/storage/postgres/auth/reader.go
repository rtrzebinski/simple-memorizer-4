package auth

import "context"

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) SignIn(ctx context.Context, email, password string) (name, userID string, err error) {
	return "name", "userID", nil
}
