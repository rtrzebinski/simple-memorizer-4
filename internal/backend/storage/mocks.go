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

func (mock *ReaderMock) RandomExercise() (models.Exercise, error) {
	return mock.Called().Get(0).(models.Exercise), nil
}

type WriterMock struct{ mock.Mock }

func NewWriterMock() *WriterMock {
	return &WriterMock{}
}

func (mock *WriterMock) StoreExercise(exercise models.Exercise) error {
	mock.Called(exercise)

	return nil
}

func (mock *WriterMock) DeleteExercise(exercise models.Exercise) error {
	mock.Called(exercise)

	return nil
}

func (mock *WriterMock) StoreLesson(lesson models.Lesson) error {
	mock.Called(lesson)

	return nil
}

func (mock *WriterMock) DeleteLesson(lesson models.Lesson) error {
	mock.Called(lesson)

	return nil
}

func (mock *WriterMock) IncrementBadAnswers(exercise models.Exercise) error {
	mock.Called(exercise)

	return nil
}

func (mock *WriterMock) IncrementGoodAnswers(exercise models.Exercise) error {
	mock.Called(exercise)

	return nil
}
