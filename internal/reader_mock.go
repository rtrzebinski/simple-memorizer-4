package internal

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
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
