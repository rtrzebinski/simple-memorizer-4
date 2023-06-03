package storage

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/mock"
)

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
