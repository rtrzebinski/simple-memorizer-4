package backend

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchLessons(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	expectedLessons := Lessons{}
	readerMock.On("FetchLessons").Return(expectedLessons, nil)

	lessons, err := service.FetchLessons()

	assert.NoError(t, err)
	assert.Equal(t, expectedLessons, lessons)
	readerMock.AssertCalled(t, "FetchLessons")
}

func TestService_HydrateLesson(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	lesson := &Lesson{}
	readerMock.On("HydrateLesson", lesson).Return(nil)

	err := service.HydrateLesson(lesson)

	assert.NoError(t, err)
	readerMock.AssertCalled(t, "HydrateLesson", lesson)
}

func TestService_FetchExercises(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	lesson := Lesson{}
	expectedExercises := Exercises{}
	readerMock.On("FetchExercises", lesson).Return(expectedExercises, nil)

	exercises, err := service.FetchExercises(lesson)

	assert.NoError(t, err)
	assert.Equal(t, expectedExercises, exercises)
	readerMock.AssertCalled(t, "FetchExercises", lesson)
}

func TestService_UpsertLesson(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	lesson := &Lesson{}
	writerMock.On("UpsertLesson", lesson).Return(nil)

	err := service.UpsertLesson(lesson)

	assert.NoError(t, err)
	writerMock.AssertCalled(t, "UpsertLesson", lesson)
}

func TestService_DeleteLesson(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	lesson := Lesson{}
	writerMock.On("DeleteLesson", lesson).Return(nil)

	err := service.DeleteLesson(lesson)

	assert.NoError(t, err)
	writerMock.AssertCalled(t, "DeleteLesson", lesson)
}

func TestService_UpsertExercise(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	exercise := &Exercise{}
	writerMock.On("UpsertExercise", exercise).Return(nil)

	err := service.UpsertExercise(exercise)

	assert.NoError(t, err)
	writerMock.AssertCalled(t, "UpsertExercise", exercise)
}

func TestService_StoreExercises(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	exercises := Exercises{}
	writerMock.On("StoreExercises", exercises).Return(nil)

	err := service.StoreExercises(exercises)

	assert.NoError(t, err)
	writerMock.AssertCalled(t, "StoreExercises", exercises)
}

func TestService_DeleteExercise(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	exercise := Exercise{}
	writerMock.On("DeleteExercise", exercise).Return(nil)

	err := service.DeleteExercise(exercise)

	assert.NoError(t, err)
	writerMock.AssertCalled(t, "DeleteExercise", exercise)
}

func TestService_PublishGoodAnswer(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	ctx := context.TODO()
	exerciseID := 123
	publisherMock.On("PublishGoodAnswer", ctx, exerciseID).Return(nil)

	err := service.PublishGoodAnswer(ctx, exerciseID)

	assert.NoError(t, err)
	publisherMock.AssertCalled(t, "PublishGoodAnswer", ctx, exerciseID)
}

func TestService_PublishBadAnswer(t *testing.T) {
	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	service := NewService(readerMock, writerMock, publisherMock)

	ctx := context.TODO()
	exerciseID := 123
	publisherMock.On("PublishBadAnswer", ctx, exerciseID).Return(nil)

	err := service.PublishBadAnswer(ctx, exerciseID)

	assert.NoError(t, err)
	publisherMock.AssertCalled(t, "PublishBadAnswer", ctx, exerciseID)
}
