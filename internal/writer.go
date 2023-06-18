package internal

import "github.com/rtrzebinski/simple-memorizer-4/internal/models"

type Writer interface {
	UpsertLesson(*models.Lesson) error
	DeleteLesson(models.Lesson) error
	UpsertExercise(*models.Exercise) error
	DeleteExercise(models.Exercise) error
	StoreResult(*models.Result) error
}
