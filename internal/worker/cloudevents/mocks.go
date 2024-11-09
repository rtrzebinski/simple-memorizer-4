package cloudevents

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) ProcessGoodAnswer(ctx context.Context, exerciseID int) error {
	args := m.Called(ctx, exerciseID)
	return args.Error(0)
}

func (m *ServiceMock) ProcessBadAnswer(ctx context.Context, exerciseID int) error {
	args := m.Called(ctx, exerciseID)
	return args.Error(0)
}
