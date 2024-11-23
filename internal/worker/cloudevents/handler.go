package cloudevents

import (
	"context"
	"fmt"

	cprotobuf "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/events"
)

const (
	GoodAnswerType = "exercise_good_answer"
	BadAnswerType  = "exercise_bad_answer"
)

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) Handle(ctx context.Context, e event.Event) error {
	switch e.Type() {
	case GoodAnswerType:
		return h.handleGoodAnswer(ctx, e)
	case BadAnswerType:
		return h.handleBadAnswer(ctx, e)
	default:
		return fmt.Errorf("event type not accepted: %s", e.Type())
	}
}

func (h *Handler) handleGoodAnswer(ctx context.Context, e event.Event) error {
	var message events.GoodAnswer

	err := cprotobuf.DecodeData(ctx, e.Data(), &message)
	if err != nil {
		return err
	}

	err = h.s.ProcessGoodAnswer(ctx, int(message.ExerciseID))
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) handleBadAnswer(ctx context.Context, e event.Event) error {
	var message events.BadAnswer

	err := cprotobuf.DecodeData(ctx, e.Data(), &message)
	if err != nil {
		return err
	}

	err = h.s.ProcessBadAnswer(ctx, int(message.ExerciseID))
	if err != nil {
		return err
	}

	return nil
}
