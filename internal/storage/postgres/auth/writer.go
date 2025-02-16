package auth

import "context"

type Writer struct {
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Register(ctx context.Context, name, email, password string) (userID string, err error) {
	return "userID", nil
}
