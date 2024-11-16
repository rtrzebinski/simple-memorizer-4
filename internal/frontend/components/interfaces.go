package components

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
)

type APIClient interface {
	FetchLessons() ([]frontend.Lesson, error)
	HydrateLesson(*frontend.Lesson) error
	FetchExercises(frontend.Lesson) ([]frontend.Exercise, error)
	UpsertLesson(frontend.Lesson) error
	DeleteLesson(frontend.Lesson) error
	UpsertExercise(frontend.Exercise) error
	StoreExercises([]frontend.Exercise) error
	DeleteExercise(frontend.Exercise) error
	StoreResult(frontend.Result) error
}
