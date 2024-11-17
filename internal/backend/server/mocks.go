package server

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() *ServiceMock {
	return &ServiceMock{}
}

func (m *ServiceMock) FetchLessons() (backend.Lessons, error) {
	args := m.Called()
	return args.Get(0).(backend.Lessons), args.Error(1)
}

func (m *ServiceMock) HydrateLesson(lesson *backend.Lesson) error {
	args := m.Called(lesson)
	return args.Error(0)
}

func (m *ServiceMock) FetchExercises(lesson backend.Lesson) (backend.Exercises, error) {
	args := m.Called(lesson)
	return args.Get(0).(backend.Exercises), args.Error(1)
}

func (m *ServiceMock) UpsertLesson(lesson *backend.Lesson) error {
	args := m.Called(lesson)
	return args.Error(0)
}

func (m *ServiceMock) DeleteLesson(lesson backend.Lesson) error {
	args := m.Called(lesson)
	return args.Error(0)
}

func (m *ServiceMock) UpsertExercise(exercise *backend.Exercise) error {
	args := m.Called(exercise)
	return args.Error(0)
}

func (m *ServiceMock) StoreExercises(exercises backend.Exercises) error {
	args := m.Called(exercises)
	return args.Error(0)
}

func (m *ServiceMock) DeleteExercise(exercise backend.Exercise) error {
	args := m.Called(exercise)
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
