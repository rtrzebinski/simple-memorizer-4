package api

import (
	"github.com/stretchr/testify/mock"
)

type CallerMock struct{ mock.Mock }

func NewCallerMock() *CallerMock {
	return &CallerMock{}
}

func (mock *CallerMock) Call(method, route string, params map[string]string, reqBody []byte) ([]byte, error) {
	return mock.Called(method, route, params, reqBody).Get(0).([]byte), nil
}
