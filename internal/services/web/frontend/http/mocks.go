package http

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type CallerMock struct{ mock.Mock }

func NewCallerMock() *CallerMock {
	return &CallerMock{}
}

func (mock *CallerMock) Call(ctx context.Context, method, route string, params map[string]string, reqBody []byte, authToken string) ([]byte, error) {
	return mock.Called(ctx, method, route, params, reqBody, authToken).Get(0).([]byte), nil
}
