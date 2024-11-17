package server

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
)

type Service interface {
	FetchLessons() (backend.Lessons, error)
	HydrateLesson(*backend.Lesson) error
	FetchExercises(backend.Lesson) (backend.Exercises, error)
	UpsertLesson(*backend.Lesson) error
	DeleteLesson(backend.Lesson) error
	UpsertExercise(*backend.Exercise) error
	StoreExercises(backend.Exercises) error
	DeleteExercise(backend.Exercise) error
	PublishGoodAnswer(ctx context.Context, exerciseID int) error
	PublishBadAnswer(ctx context.Context, exerciseID int) error
}
