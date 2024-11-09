package cloudevents

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

type HandlerSuite struct {
	suite.Suite
	serviceMock *ServiceMock
	handler     *Handler
}

func (suite *HandlerSuite) SetupTest() {
	suite.serviceMock = new(ServiceMock)
	suite.handler = NewHandler(suite.serviceMock)
}

func (suite *HandlerSuite) TestHandler_Handle_GoodAnswer_Success() {
	ctx := context.Background()
	exerciseID := 1
	message := events.GoodAnswer{ExerciseID: uint32(exerciseID)}
	e := event.New()
	e.SetType(GoodAnswerType)
	err := e.SetData(cloudevents.ApplicationJSON, &message)
	suite.NoError(err)

	suite.serviceMock.On("ProcessGoodAnswer", ctx, exerciseID).Return(nil)

	err = suite.handler.Handle(ctx, e)

	suite.NoError(err)
	suite.serviceMock.AssertCalled(suite.T(), "ProcessGoodAnswer", ctx, exerciseID)
}

func (suite *HandlerSuite) TestHandler_Handle_GoodAnswer_Error() {
	ctx := context.Background()
	exerciseID := 1
	message := events.GoodAnswer{ExerciseID: uint32(exerciseID)}
	e := event.New()
	e.SetType(GoodAnswerType)
	err := e.SetData(cloudevents.ApplicationJSON, &message)
	suite.NoError(err)

	suite.serviceMock.On("ProcessGoodAnswer", ctx, exerciseID).Return(assert.AnError)

	err = suite.handler.Handle(ctx, e)

	suite.Error(err)
	suite.Contains(err.Error(), assert.AnError.Error())
	suite.serviceMock.AssertCalled(suite.T(), "ProcessGoodAnswer", ctx, exerciseID)
}

func (suite *HandlerSuite) TestHandler_Handle_BadAnswer_Success() {
	ctx := context.Background()
	exerciseID := 2
	message := events.BadAnswer{ExerciseID: uint32(exerciseID)}
	e := event.New()
	e.SetType(BadAnswerType)
	err := e.SetData(cloudevents.ApplicationJSON, &message)
	suite.NoError(err)

	suite.serviceMock.On("ProcessBadAnswer", ctx, exerciseID).Return(nil)

	err = suite.handler.Handle(ctx, e)

	suite.NoError(err)
	suite.serviceMock.AssertCalled(suite.T(), "ProcessBadAnswer", ctx, exerciseID)
}

func (suite *HandlerSuite) TestHandler_Handle_BadAnswer_Error() {
	ctx := context.Background()
	exerciseID := 2
	message := events.BadAnswer{ExerciseID: uint32(exerciseID)}
	e := event.New()
	e.SetType(BadAnswerType)
	err := e.SetData(cloudevents.ApplicationJSON, &message)
	suite.NoError(err)

	suite.serviceMock.On("ProcessBadAnswer", ctx, exerciseID).Return(assert.AnError)

	err = suite.handler.Handle(ctx, e)

	suite.Error(err)
	suite.Contains(err.Error(), assert.AnError.Error())
	suite.serviceMock.AssertCalled(suite.T(), "ProcessBadAnswer", ctx, exerciseID)
}

func (suite *HandlerSuite) TestHandler_Handle_UnknownEventType() {
	ctx := context.Background()
	e := event.New()
	e.SetType("unknown_event_type")

	err := suite.handler.Handle(ctx, e)

	suite.Error(err)
	suite.Contains(err.Error(), "event type not accepted: unknown_event_type")
}
