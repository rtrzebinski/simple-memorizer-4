package cloudevents

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/mock"
)

type SenderMock struct {
	mock.Mock
}

func (m *SenderMock) Send(ctx context.Context, e event.Event) cloudevents.Result {
	args := m.Called(ctx, e)
	return args.Get(0).(cloudevents.Result)
}
