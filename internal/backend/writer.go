package backend

import "github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"

type Writer interface {
	UpsertLesson(*models.Lesson) error
	DeleteLesson(models.Lesson) error
	UpsertExercise(*models.Exercise) error
	StoreExercises(models.Exercises) error
	DeleteExercise(models.Exercise) error
	StoreResult(*models.Result) error
}
