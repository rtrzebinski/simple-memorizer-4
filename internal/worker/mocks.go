package worker

import (
	"github.com/stretchr/testify/mock"
)

type WriterMock struct {
	mock.Mock
}

func (m *WriterMock) StoreResult(result *Result) error {
	args := m.Called(result)
	return args.Error(0)
}
