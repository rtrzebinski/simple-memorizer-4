package storage

import "github.com/rtrzebinski/simple-memorizer-4/internal/models"

type Reader interface {
	ExercisesOfLesson(lessonId int) (models.Exercises, error)
	RandomExerciseOfLesson(lessonId int) (models.Exercise, error)
	AllLessons() (models.Lessons, error)
}
