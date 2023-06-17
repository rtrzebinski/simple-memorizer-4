package internal

import "github.com/rtrzebinski/simple-memorizer-4/internal/models"

type Writer interface {
	StoreLesson(*models.Lesson) error
	DeleteLesson(models.Lesson) error
	StoreExercise(*models.Exercise) error
	DeleteExercise(models.Exercise) error
	StoreAnswer(*models.Answer) error
}
