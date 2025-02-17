package http

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type CallerMock struct{ mock.Mock }

func NewCallerMock() *CallerMock {
	return &CallerMock{}
}

func (mock *CallerMock) Call(ctx context.Context, method, route string, params map[string]string, reqBody []byte) ([]byte, error) {
	return mock.Called(ctx, method, route, params, reqBody).Get(0).([]byte), nil
}
