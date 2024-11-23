package cloudevents

import (
	"context"
	"fmt"
	"log/slog"

	cprotobuf "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/events"
	"google.golang.org/protobuf/proto"
)

const (
	ContentType    = "application/protobuf"
	GoodAnswerType = "exercise_good_answer"
	BadAnswerType  = "exercise_bad_answer"
	Source         = "github.com/rtrzebinski/simple-memorizer-4"
)

type Publisher struct {
	s Sender
}

func NewPublisher(s Sender) *Publisher {
	return &Publisher{
		s: s,
	}
}

func (p *Publisher) PublishGoodAnswer(ctx context.Context, exerciseID int) error {
	message := events.GoodAnswer{
		ExerciseID: uint32(exerciseID),
	}

	err := p.publish(ctx, &message, GoodAnswerType)
	if err != nil {
		return fmt.Errorf("publish good answer: %w", err)
	}

	slog.Info("published events.GoodAnswer", slog.Int("exerciseID", exerciseID), slog.String("service", "web"))

	return nil
}

func (p *Publisher) PublishBadAnswer(ctx context.Context, exerciseID int) error {
	message := events.BadAnswer{
		ExerciseID: uint32(exerciseID),
	}

	err := p.publish(ctx, &message, BadAnswerType)
	if err != nil {
		return fmt.Errorf("publish bad answer: %w", err)
	}

	slog.Info("published events.BadAnswer", slog.Int("exerciseID", exerciseID), slog.String("service", "web"))

	return nil
}

func (p *Publisher) publish(ctx context.Context, message proto.Message, eType string) error {
	protobufData, err := cprotobuf.EncodeData(ctx, message)
	if err != nil {
		return fmt.Errorf("encode data: %w", err)
	}

	ce := cloudevents.NewEvent()
	ce.SetType(eType)
	ce.SetSource(Source)

	if err := ce.SetData(ContentType, protobufData); err != nil {
		return fmt.Errorf("set data: %w", err)
	}

	if result := p.s.Send(ctx, ce); cloudevents.IsUndelivered(result) {
		return fmt.Errorf("send event: %w", result)
	}

	return nil
}
