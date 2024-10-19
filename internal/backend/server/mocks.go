package server

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/stretchr/testify/mock"
)

type ReaderMock struct{ mock.Mock }

func NewReaderMock() *ReaderMock {
	return &ReaderMock{}
}

func (mock *ReaderMock) FetchLessons() (models.Lessons, error) {
	return mock.Called().Get(0).(models.Lessons), nil
}

func (mock *ReaderMock) HydrateLesson(lesson *models.Lesson) error {
	mock.Called(lesson)

	return nil
}

func (mock *ReaderMock) FetchExercises(lesson models.Lesson) (models.Exercises, error) {
	return mock.Called(lesson).Get(0).(models.Exercises), nil
}

type WriterMock struct{ mock.Mock }

func NewWriterMock() *WriterMock {
	return &WriterMock{}
}

func (mock *WriterMock) UpsertLesson(lesson *models.Lesson) error {
	mock.Called(lesson)

	return nil
}

func (mock *WriterMock) DeleteLesson(lesson models.Lesson) error {
	mock.Called(lesson)

	return nil
}

func (mock *WriterMock) UpsertExercise(exercise *models.Exercise) error {
	mock.Called(exercise)

	return nil
}

func (mock *WriterMock) StoreExercises(exercises models.Exercises) error {
	mock.Called(exercises)

	return nil
}

func (mock *WriterMock) DeleteExercise(exercise models.Exercise) error {
	mock.Called(exercise)

	return nil
}

func (mock *WriterMock) StoreResult(result *models.Result) error {
	mock.Called(result)

	return nil
}
