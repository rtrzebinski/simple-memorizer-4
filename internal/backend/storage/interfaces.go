package storage

import "github.com/rtrzebinski/simple-memorizer-4/internal/models"

type Reader interface {
	Exercises() (models.Exercises, error)
	RandomExercise() (models.Exercise, error)
}

type Writer interface {
	StoreExercise(models.Exercise) error
	IncrementBadAnswers(models.Exercise) error
	IncrementGoodAnswers(models.Exercise) error
}
