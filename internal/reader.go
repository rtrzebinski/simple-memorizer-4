package internal

import "github.com/rtrzebinski/simple-memorizer-4/internal/models"

type Reader interface {
	FetchAllLessons() (models.Lessons, error)
	HydrateLesson(*models.Lesson) error
	FetchExercisesOfLesson(models.Lesson) (models.Exercises, error)
}
