package components

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
)

type APIClient interface {
	FetchLessons(ctx context.Context) ([]frontend.Lesson, error)
	HydrateLesson(ctx context.Context, lesson *frontend.Lesson) error
	FetchExercises(ctx context.Context, lesson frontend.Lesson) ([]frontend.Exercise, error)
	UpsertLesson(ctx context.Context, lesson frontend.Lesson) error
	DeleteLesson(ctx context.Context, lesson frontend.Lesson) error
	UpsertExercise(ctx context.Context, exercise frontend.Exercise) error
	StoreExercises(ctx context.Context, exercises []frontend.Exercise) error
	DeleteExercise(ctx context.Context, exercise frontend.Exercise) error
	StoreResult(ctx context.Context, result frontend.Result) error
}
