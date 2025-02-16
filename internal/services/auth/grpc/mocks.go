package grpc

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) Register(ctx context.Context, name, email, password string) (authToken string, err error) {
	called := m.Called(ctx, name, email, password)
	return called.Get(0).(string), called.Error(1)
}

func (m *ServiceMock) SignIn(ctx context.Context, email, password string) (authToken string, err error) {
	called := m.Called(ctx, email, password)
	return called.Get(0).(string), called.Error(1)
}
