package http

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type DoerMock struct{ mock.Mock }

func (c *DoerMock) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}
