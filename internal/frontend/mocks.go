package frontend

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type HttpClientMock struct{ mock.Mock }

func (c *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}
