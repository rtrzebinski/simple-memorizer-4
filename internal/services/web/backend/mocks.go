package backend

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ReaderMock struct{ mock.Mock }

func NewReaderMock() *ReaderMock {
	return &ReaderMock{}
}

func (mock *ReaderMock) FetchLessons(ctx context.Context) (Lessons, error) {
	return mock.Called(ctx).Get(0).(Lessons), nil
}

func (mock *ReaderMock) HydrateLesson(ctx context.Context, lesson *Lesson) error {
	mock.Called(ctx, lesson)

	return nil
}

func (mock *ReaderMock) FetchExercises(ctx context.Context, lesson Lesson) (Exercises, error) {
	return mock.Called(ctx, lesson).Get(0).(Exercises), nil
}

type WriterMock struct{ mock.Mock }

func NewWriterMock() *WriterMock {
	return &WriterMock{}
}

func (mock *WriterMock) UpsertLesson(ctx context.Context, lesson *Lesson) error {
	mock.Called(ctx, lesson)

	return nil
}

func (mock *WriterMock) DeleteLesson(ctx context.Context, lesson Lesson) error {
	mock.Called(ctx, lesson)

	return nil
}

func (mock *WriterMock) UpsertExercise(ctx context.Context, exercise *Exercise) error {
	mock.Called(ctx, exercise)

	return nil
}

func (mock *WriterMock) StoreExercises(ctx context.Context, exercises Exercises) error {
	mock.Called(ctx, exercises)

	return nil
}

func (mock *WriterMock) DeleteExercise(ctx context.Context, exercise Exercise) error {
	mock.Called(ctx, exercise)

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

type AuthClientMock struct{ mock.Mock }

func NewAuthClientMock() *AuthClientMock {
	return &AuthClientMock{}
}

func (mock *AuthClientMock) Register(ctx context.Context, name, email, password string) (accessToken string, err error) {
	args := mock.Called(ctx, name, email, password)

	return args.String(0), args.Error(1)
}

func (mock *AuthClientMock) SignIn(ctx context.Context, email, password string) (accessToken string, err error) {
	args := mock.Called(ctx, email, password)

	return args.String(0), args.Error(1)
}
