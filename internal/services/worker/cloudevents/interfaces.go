package cloudevents

import (
	"context"
)

type Service interface {
	ProcessGoodAnswer(ctx context.Context, exerciseID int) error
	ProcessBadAnswer(ctx context.Context, exerciseID int) error
}
