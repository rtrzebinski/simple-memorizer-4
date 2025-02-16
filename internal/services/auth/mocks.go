package auth

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ReaderMock struct {
	mock.Mock
}

func (m *ReaderMock) FetchUser(ctx context.Context, email string) (string, string, string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

type WriterMock struct {
	mock.Mock
}

func (m *WriterMock) StoreUser(ctx context.Context, name string, email string, password string) (string, error) {
	args := m.Called(ctx, name, email, password)
	return args.String(0), args.Error(1)
}
