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

func (m *ServiceMock) HydrateLesson(ctx context.Context, userID string, lesson *backend.Lesson) error {
	args := m.Called(ctx, userID, lesson)
	return args.Error(0)
}

func (m *ServiceMock) FetchExercises(ctx context.Context, userID string, lesson backend.Lesson, oldestExerciseID int) (backend.Exercises, error) {
	args := m.Called(ctx, userID, lesson, oldestExerciseID)
	return args.Get(0).(backend.Exercises), args.Error(1)
}

func (m *ServiceMock) UpsertLesson(ctx context.Context, userID string, lesson *backend.Lesson) error {
	args := m.Called(ctx, userID, lesson)
	return args.Error(0)
}

func (m *ServiceMock) DeleteLesson(ctx context.Context, userID string, lesson backend.Lesson) error {
	args := m.Called(ctx, userID, lesson)
	return args.Error(0)
}

func (m *ServiceMock) UpsertExercise(ctx context.Context, userID string, exercise *backend.Exercise) error {
	args := m.Called(ctx, userID, exercise)
	return args.Error(0)
}

func (m *ServiceMock) StoreExercises(ctx context.Context, userID string, exercises backend.Exercises) error {
	args := m.Called(ctx, userID, exercises)
	return args.Error(0)
}

func (m *ServiceMock) DeleteExercise(ctx context.Context, userID string, exercise backend.Exercise) error {
	args := m.Called(ctx, userID, exercise)
	return args.Error(0)
}

func (m *ServiceMock) PublishGoodAnswer(ctx context.Context, userID string, exerciseID int) error {
	args := m.Called(ctx, userID, exerciseID)
	return args.Error(0)
}

func (m *ServiceMock) PublishBadAnswer(ctx context.Context, userID string, exerciseID int) error {
	args := m.Called(ctx, userID, exerciseID)
	return args.Error(0)
}

func (m *ServiceMock) Register(ctx context.Context, firstName, lastName, email, password string) (backend.Tokens, error) {
	args := m.Called(ctx, firstName, lastName, email, password)
	return args.Get(0).(backend.Tokens), args.Error(1)
}

func (m *ServiceMock) SignIn(ctx context.Context, email, password string) (backend.Tokens, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(backend.Tokens), args.Error(1)
}

func (m *ServiceMock) Revoke(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

type TokenVerifierMock struct {
	mock.Mock
}

func NewTokenVerifierMock() *TokenVerifierMock {
	return &TokenVerifierMock{}
}

func (m *TokenVerifierMock) VerifyAndUser(ctx context.Context, accessToken string) (*backend.User, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(*backend.User), args.Error(1)
}

type TokenRefresherMock struct {
	mock.Mock
}

func NewTokenRefresherMock() *TokenRefresherMock {
	return &TokenRefresherMock{}
}

func (m *TokenRefresherMock) Refresh(ctx context.Context, refreshToken string) (backend.Tokens, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(backend.Tokens), args.Error(1)
}
