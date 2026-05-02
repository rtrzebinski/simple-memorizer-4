package pubsub

import (
	"context"
)

type Service interface {
	ProcessGoodAnswer(ctx context.Context, userID string, exerciseID int) error
	ProcessBadAnswer(ctx context.Context, userID string, exerciseID int) error
}
