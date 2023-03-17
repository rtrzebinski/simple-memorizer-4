package storage

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/mock"
)

type ReaderMock struct{ mock.Mock }

func NewReaderMock() *ReaderMock {
	return &ReaderMock{}
}

func (mock *ReaderMock) RandomExercise() models.Exercise {
	return mock.Called().Get(0).(models.Exercise)
}

type WriterMock struct{ mock.Mock }

func NewWriterMock() *WriterMock {
	return &WriterMock{}
}

func (mock *WriterMock) IncrementBadAnswers(exerciseId int) {
	mock.Called(exerciseId)
}

func (mock *WriterMock) IncrementGoodAnswers(exerciseId int) {
	mock.Called(exerciseId)
}
