package server

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/stretchr/testify/mock"
)

type ReaderMock struct{ mock.Mock }

func NewReaderMock() *ReaderMock {
	return &ReaderMock{}
}

func (mock *ReaderMock) FetchLessons() (backend.Lessons, error) {
	return mock.Called().Get(0).(backend.Lessons), nil
}

func (mock *ReaderMock) HydrateLesson(lesson *backend.Lesson) error {
	mock.Called(lesson)

	return nil
}

func (mock *ReaderMock) FetchExercises(lesson backend.Lesson) (backend.Exercises, error) {
	return mock.Called(lesson).Get(0).(backend.Exercises), nil
}

type WriterMock struct{ mock.Mock }

func NewWriterMock() *WriterMock {
	return &WriterMock{}
}

func (mock *WriterMock) UpsertLesson(lesson *backend.Lesson) error {
	mock.Called(lesson)

	return nil
}

func (mock *WriterMock) DeleteLesson(lesson backend.Lesson) error {
	mock.Called(lesson)

	return nil
}

func (mock *WriterMock) UpsertExercise(exercise *backend.Exercise) error {
	mock.Called(exercise)

	return nil
}

func (mock *WriterMock) StoreExercises(exercises backend.Exercises) error {
	mock.Called(exercises)

	return nil
}

func (mock *WriterMock) DeleteExercise(exercise backend.Exercise) error {
	mock.Called(exercise)

	return nil
}

type PublisherMock struct{ mock.Mock }

func NewPublisherMock() *PublisherMock {
	return &PublisherMock{}
}

func (mock *PublisherMock) PublishGoodAnswer(ctx context.Context, exerciseId int) error {
	mock.Called(ctx, exerciseId)

	return nil
}

func (mock *PublisherMock) PublishBadAnswer(ctx context.Context, exerciseId int) error {
	mock.Called(ctx, exerciseId)

	return nil
}
