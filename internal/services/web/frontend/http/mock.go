package http

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/stretchr/testify/mock"
)

type CallerMock struct{ mock.Mock }

func NewCallerMock() *CallerMock {
	return &CallerMock{}
}

func (mock *CallerMock) Call(ctx app.Context, method, route string, params map[string]string, reqBody []byte) ([]byte, error) {
	return mock.Called(ctx, method, route, params, reqBody).Get(0).([]byte), nil
}
