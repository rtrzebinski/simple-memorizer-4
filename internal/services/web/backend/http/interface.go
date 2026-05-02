package http

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type Service interface {
	FetchLessons(ctx context.Context, userID string) (lessons backend.Lessons, err error)
	HydrateLesson(ctx context.Context, lesson *backend.Lesson, userID string) error
	FetchExercises(ctx context.Context, lesson backend.Lesson, oldestExerciseID int, userID string) (backend.Exercises, error)
	UpsertLesson(ctx context.Context, lesson *backend.Lesson, userID string) error
	DeleteLesson(ctx context.Context, lesson backend.Lesson, userID string) error
	UpsertExercise(ctx context.Context, exercise *backend.Exercise, userID string) error
	StoreExercises(ctx context.Context, exercise backend.Exercises, userID string) error
	DeleteExercise(ctx context.Context, exercise backend.Exercise, userID string) error
	PublishGoodAnswer(ctx context.Context, exerciseID int, userID string) error
	PublishBadAnswer(ctx context.Context, exerciseID int, userID string) error
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
