package http

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type Service interface {
	FetchLessons(ctx context.Context, userID string) (lessons backend.Lessons, err error)
	HydrateLesson(ctx context.Context, userID string, lesson *backend.Lesson) error
	FetchExercises(ctx context.Context, userID string, lesson backend.Lesson, oldestExerciseID int) (backend.Exercises, error)
	UpsertLesson(ctx context.Context, userID string, lesson *backend.Lesson) error
	DeleteLesson(ctx context.Context, userID string, lesson backend.Lesson) error
	UpsertExercise(ctx context.Context, userID string, exercise *backend.Exercise) error
	StoreExercises(ctx context.Context, userID string, exercise backend.Exercises) error
	DeleteExercise(ctx context.Context, userID string, exercise backend.Exercise) error
	PublishGoodAnswer(ctx context.Context, userID string, exerciseID int) error
	PublishBadAnswer(ctx context.Context, userID string, exerciseID int) error
	Register(ctx context.Context, firstName, lastName, email, password string) (backend.Tokens, error)
	SignIn(ctx context.Context, email, password string) (backend.Tokens, error)
	Revoke(ctx context.Context, refreshToken string) error
}

type TokenRefresher interface {
	Refresh(ctx context.Context, refreshToken string) (backend.Tokens, error)
}

type TokenVerifier interface {
	VerifyAndUser(ctx context.Context, accessToken string) (*backend.User, error)
}
