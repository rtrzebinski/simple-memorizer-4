package http

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type Service interface {
	FetchLessons(ctx context.Context) (lessons backend.Lessons, err error)
	HydrateLesson(ctx context.Context, lesson *backend.Lesson) error
	FetchExercises(ctx context.Context, lesson backend.Lesson) (backend.Exercises, error)
	UpsertLesson(ctx context.Context, lesson *backend.Lesson) error
	DeleteLesson(ctx context.Context, lesson backend.Lesson) error
	UpsertExercise(ctx context.Context, exercise *backend.Exercise) error
	StoreExercises(ctx context.Context, exercise backend.Exercises) error
	DeleteExercise(ctx context.Context, exercise backend.Exercise) error
	PublishGoodAnswer(ctx context.Context, exerciseID int) error
	PublishBadAnswer(ctx context.Context, exerciseID int) error
}
