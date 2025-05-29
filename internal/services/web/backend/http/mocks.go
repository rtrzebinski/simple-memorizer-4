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

func (m *ServiceMock) FetchLessons(ctx context.Context, userID string) (backend.Lessons, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(backend.Lessons), args.Error(1)
}

func (m *ServiceMock) HydrateLesson(ctx context.Context, lesson *backend.Lesson, userID string) error {
	args := m.Called(ctx, lesson, userID)
	return args.Error(0)
}

func (m *ServiceMock) FetchExercises(ctx context.Context, lesson backend.Lesson, oldestExerciseID int, userID string) (backend.Exercises, error) {
	args := m.Called(ctx, lesson, oldestExerciseID, userID)
	return args.Get(0).(backend.Exercises), args.Error(1)
}

func (m *ServiceMock) UpsertLesson(ctx context.Context, lesson *backend.Lesson, userID string) error {
	args := m.Called(ctx, lesson, userID)
	return args.Error(0)
}

func (m *ServiceMock) DeleteLesson(ctx context.Context, lesson backend.Lesson, userID string) error {
	args := m.Called(ctx, lesson, userID)
	return args.Error(0)
}

func (m *ServiceMock) UpsertExercise(ctx context.Context, exercise *backend.Exercise, userID string) error {
	args := m.Called(ctx, exercise, userID)
	return args.Error(0)
}

func (m *ServiceMock) StoreExercises(ctx context.Context, exercises backend.Exercises, userID string) error {
	args := m.Called(ctx, exercises, userID)
	return args.Error(0)
}

func (m *ServiceMock) DeleteExercise(ctx context.Context, exercise backend.Exercise, userID string) error {
	args := m.Called(ctx, exercise, userID)
	return args.Error(0)
}

func (m *ServiceMock) PublishGoodAnswer(ctx context.Context, exerciseID int, userID string) error {
	args := m.Called(ctx, exerciseID, userID)
	return args.Error(0)
}

func (m *ServiceMock) PublishBadAnswer(ctx context.Context, exerciseID int, userID string) error {
	args := m.Called(ctx, exerciseID, userID)
	return args.Error(0)
}

func (m *ServiceMock) Register(ctx context.Context, name, email, password string) (accessToken string, err error) {
	args := m.Called(ctx, name, email, password)
	return args.String(0), args.Error(1)
}

func (m *ServiceMock) SignIn(ctx context.Context, email, password string) (accessToken string, err error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}
