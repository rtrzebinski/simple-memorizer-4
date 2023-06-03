package storage

import "github.com/rtrzebinski/simple-memorizer-4/internal/models"

type Reader interface {
	FetchAllLessons() (models.Lessons, error)
	FetchExercisesOfLesson(models.Lesson) (models.Exercises, error)
	FetchRandomExerciseOfLesson(models.Lesson) (models.Exercise, error)
}
