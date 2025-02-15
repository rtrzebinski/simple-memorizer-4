package cloudevents

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol"
)

type Sender interface {
	Send(ctx context.Context, e event.Event) protocol.Result
}
