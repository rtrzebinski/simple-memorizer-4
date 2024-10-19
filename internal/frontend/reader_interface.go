package frontend

import "github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"

type Reader interface {
	FetchLessons() (models.Lessons, error)
	HydrateLesson(*models.Lesson) error
	FetchExercises(models.Lesson) (models.Exercises, error)
}
