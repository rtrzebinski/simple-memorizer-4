package http

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type Service interface {
	FetchLessons(ctx context.Context, userID string) (lessons backend.Lessons, err error)
	HydrateLesson(ctx context.Context, lesson *backend.Lesson, userID string) error
	FetchExercises(ctx context.Context, lesson backend.Lesson, userID string) (backend.Exercises, error)
	UpsertLesson(ctx context.Context, lesson *backend.Lesson, userID string) error
	DeleteLesson(ctx context.Context, lesson backend.Lesson, userID string) error
	UpsertExercise(ctx context.Context, exercise *backend.Exercise, userID string) error
	StoreExercises(ctx context.Context, exercise backend.Exercises, userID string) error
	DeleteExercise(ctx context.Context, exercise backend.Exercise, userID string) error
	PublishGoodAnswer(ctx context.Context, exerciseID int, userID string) error
	PublishBadAnswer(ctx context.Context, exerciseID int, userID string) error
	Register(ctx context.Context, name, email, password string) (accessToken string, err error)
	SignIn(ctx context.Context, email, password string) (accessToken string, err error)
}
