package storage

import "github.com/rtrzebinski/simple-memorizer-4/internal/models"

type Reader interface {
	ExercisesOfLesson(models.Lesson) (models.Exercises, error)
	RandomExerciseOfLesson(models.Lesson) (models.Exercise, error)
	AllLessons() (models.Lessons, error)
}
