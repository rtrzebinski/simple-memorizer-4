package backend

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchLessons(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	expectedLessons := Lessons{}
	readerMock.On("FetchLessons", ctx).Return(expectedLessons, nil)

	lessons, err := service.FetchLessons(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedLessons, lessons)
	readerMock.AssertExpectations(t)
}

func TestService_HydrateLesson(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	lesson := &Lesson{}
	readerMock.On("HydrateLesson", ctx, lesson).Return(nil)

	err := service.HydrateLesson(ctx, lesson)

	assert.NoError(t, err)
	readerMock.AssertExpectations(t)
}

func TestService_FetchExercises(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	lesson := Lesson{}
	expectedExercises := Exercises{}
	readerMock.On("FetchExercises", ctx, lesson).Return(expectedExercises, nil)

	exercises, err := service.FetchExercises(ctx, lesson)

	assert.NoError(t, err)
	assert.Equal(t, expectedExercises, exercises)
	readerMock.AssertExpectations(t)
}

func TestService_UpsertLesson(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	lesson := &Lesson{}
	writerMock.On("UpsertLesson", ctx, lesson).Return(nil)

	err := service.UpsertLesson(ctx, lesson)

	assert.NoError(t, err)
	writerMock.AssertExpectations(t)
}

func TestService_DeleteLesson(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	lesson := Lesson{}
	writerMock.On("DeleteLesson", ctx, lesson).Return(nil)

	err := service.DeleteLesson(ctx, lesson)

	assert.NoError(t, err)
	writerMock.AssertExpectations(t)
}

func TestService_UpsertExercise(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	exercise := &Exercise{}
	writerMock.On("UpsertExercise", ctx, exercise).Return(nil)

	err := service.UpsertExercise(ctx, exercise)

	assert.NoError(t, err)
	writerMock.AssertExpectations(t)
}

func TestService_StoreExercises(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	exercises := Exercises{}
	writerMock.On("StoreExercises", ctx, exercises).Return(nil)

	err := service.StoreExercises(ctx, exercises)

	assert.NoError(t, err)
	writerMock.AssertExpectations(t)
}

func TestService_DeleteExercise(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	exercise := Exercise{}
	writerMock.On("DeleteExercise", ctx, exercise).Return(nil)

	err := service.DeleteExercise(ctx, exercise)

	assert.NoError(t, err)
	writerMock.AssertExpectations(t)
}

func TestService_PublishGoodAnswer(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	exerciseID := 123
	publisherMock.On("PublishGoodAnswer", ctx, exerciseID).Return(nil)

	err := service.PublishGoodAnswer(ctx, exerciseID)

	assert.NoError(t, err)
	publisherMock.AssertExpectations(t)
}

func TestService_PublishBadAnswer(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	exerciseID := 123
	publisherMock.On("PublishBadAnswer", ctx, exerciseID).Return(nil)

	err := service.PublishBadAnswer(ctx, exerciseID)

	assert.NoError(t, err)
	publisherMock.AssertExpectations(t)
}

func TestService_Register(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	name := "name"
	email := "email"
	password := "password"
	accessToken := "accessToken"
	authClientMock.On("Register", ctx, name, email, password).Return(accessToken, nil)

	token, err := service.Register(ctx, name, email, password)

	assert.NoError(t, err)
	assert.Equal(t, accessToken, token)
	authClientMock.AssertExpectations(t)
}

func TestService_SignIn(t *testing.T) {
	ctx := context.Background()

	readerMock := NewReaderMock()
	writerMock := NewWriterMock()
	publisherMock := NewPublisherMock()
	authClientMock := NewAuthClientMock()
	service := NewService(readerMock, writerMock, publisherMock, authClientMock)

	email := "email"
	password := "password"
	accessToken := "accessToken"
	authClientMock.On("SignIn", ctx, email, password).Return(accessToken, nil)

	token, err := service.SignIn(ctx, email, password)

	assert.NoError(t, err)
	assert.Equal(t, accessToken, token)
	authClientMock.AssertExpectations(t)
}
