package storage

import "github.com/rtrzebinski/simple-memorizer-4/internal/models"

type Reader interface {
	ExercisesOfLesson(lessonId int) (models.Exercises, error)
	RandomExerciseOfLesson(lessonId int) (models.Exercise, error)
	AllLessons() (models.Lessons, error)
}

type Writer interface {
	DeleteExercise(models.Exercise) error
	StoreExercise(models.Exercise) error
	DeleteLesson(models.Lesson) error
	StoreLesson(models.Lesson) error
	IncrementBadAnswers(models.Exercise) error
	IncrementGoodAnswers(models.Exercise) error
}
