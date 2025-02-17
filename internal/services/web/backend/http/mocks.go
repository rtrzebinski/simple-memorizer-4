package http

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() *ServiceMock {
	return &ServiceMock{}
}

func (m *ServiceMock) FetchLessons(ctx context.Context) (backend.Lessons, error) {
	args := m.Called(ctx)
	return args.Get(0).(backend.Lessons), args.Error(1)
}

func (m *ServiceMock) HydrateLesson(ctx context.Context, lesson *backend.Lesson) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *ServiceMock) FetchExercises(ctx context.Context, lesson backend.Lesson) (backend.Exercises, error) {
	args := m.Called(ctx, lesson)
	return args.Get(0).(backend.Exercises), args.Error(1)
}

func (m *ServiceMock) UpsertLesson(ctx context.Context, lesson *backend.Lesson) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *ServiceMock) DeleteLesson(ctx context.Context, lesson backend.Lesson) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *ServiceMock) UpsertExercise(ctx context.Context, exercise *backend.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

func (m *ServiceMock) StoreExercises(ctx context.Context, exercises backend.Exercises) error {
	args := m.Called(ctx, exercises)
	return args.Error(0)
}

func (m *ServiceMock) DeleteExercise(ctx context.Context, exercise backend.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

func (m *ServiceMock) PublishGoodAnswer(ctx context.Context, exerciseID int) error {
	args := m.Called(ctx, exerciseID)
	return args.Error(0)
}

func (m *ServiceMock) PublishBadAnswer(ctx context.Context, exerciseID int) error {
	args := m.Called(ctx, exerciseID)
	return args.Error(0)
}
