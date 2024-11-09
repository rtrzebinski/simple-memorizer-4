package cloudevents

import (
	"context"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PublisherSuite struct {
	suite.Suite
	senderMock *SenderMock
	publisher  *Publisher
}

func (suite *PublisherSuite) SetupTest() {
	suite.senderMock = new(SenderMock)
	suite.publisher = NewPublisher(suite.senderMock)
}

func (suite *PublisherSuite) TestPublisher_Publish_GoodAnswer() {
	ctx := context.Background()
	exerciseID := 123

	suite.senderMock.On("Send", ctx, mock.MatchedBy(func(e event.Event) bool {
		assert.EqualValues(suite.T(), ContentType, e.DataContentType())
		assert.EqualValues(suite.T(), GoodAnswerType, e.Type())
		assert.EqualValues(suite.T(), Source, e.Source())

		// Decode and verify event data
		var goodAnswer events.GoodAnswer
		err := e.DataAs(&goodAnswer)
		suite.NoError(err)
		assert.EqualValues(suite.T(), uint32(exerciseID), goodAnswer.ExerciseID)

		return true
	})).Return(cloudevents.ResultACK)

	err := suite.publisher.PublishGoodAnswer(ctx, exerciseID)
	suite.NoError(err)

	suite.senderMock.AssertExpectations(suite.T())
}

func (suite *PublisherSuite) TestPublisher_Publish_BadAnswer() {
	ctx := context.Background()
	exerciseID := 456

	suite.senderMock.On("Send", ctx, mock.MatchedBy(func(e event.Event) bool {
		assert.EqualValues(suite.T(), ContentType, e.DataContentType())
		assert.EqualValues(suite.T(), BadAnswerType, e.Type())
		assert.EqualValues(suite.T(), Source, e.Source())

		// Decode and verify event data
		var badAnswer events.BadAnswer
		err := e.DataAs(&badAnswer)
		suite.NoError(err)
		assert.EqualValues(suite.T(), uint32(exerciseID), badAnswer.ExerciseID)

		return true
	})).Return(cloudevents.ResultACK)

	err := suite.publisher.PublishBadAnswer(ctx, exerciseID)
	suite.NoError(err)

	suite.senderMock.AssertExpectations(suite.T())
}

func TestPublisherSuite(t *testing.T) {
	suite.Run(t, new(PublisherSuite))
}
