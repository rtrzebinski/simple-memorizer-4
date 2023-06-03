package storage

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/mock"
)

type ReaderMock struct{ mock.Mock }

func NewReaderMock() *ReaderMock {
	return &ReaderMock{}
}

func (mock *ReaderMock) ExercisesOfLesson(lessonId int) (models.Exercises, error) {
	return mock.Called(lessonId).Get(0).(models.Exercises), nil
}

func (mock *ReaderMock) AllLessons() (models.Lessons, error) {
	return mock.Called().Get(0).(models.Lessons), nil
}

func (mock *ReaderMock) RandomExerciseOfLesson(lessonId int) (models.Exercise, error) {
	return mock.Called(lessonId).Get(0).(models.Exercise), nil
}
