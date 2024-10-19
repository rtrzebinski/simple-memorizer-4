package frontend

import (
	"github.com/stretchr/testify/mock"
)

type ClientMock struct{ mock.Mock }

func NewClientMock() *ClientMock {
	return &ClientMock{}
}

func (mock *ClientMock) Call(method, route string, params map[string]string, reqBody []byte) ([]byte, error) {
	return mock.Called(method, route, params, reqBody).Get(0).([]byte), nil
}
