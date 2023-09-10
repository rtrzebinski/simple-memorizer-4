package internal

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/mock"
)

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
