package storage

import "github.com/rtrzebinski/simple-memorizer-go/internal/models"

type Reader interface {
	RandomExercise() models.Exercise
}
