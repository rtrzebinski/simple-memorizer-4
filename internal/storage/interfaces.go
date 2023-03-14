package storage

import "github.com/rtrzebinski/simple-memorizer-go/internal/models"

type Reader interface {
	RandomExercise() models.Exercise
}

type Writer interface {
	IncrementGoodAnswers(exerciseId int)
	IncrementBadAnswers(exerciseId int)
}
