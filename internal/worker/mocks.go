package worker

import (
	"github.com/stretchr/testify/mock"
)

type ReaderMock struct {
	mock.Mock
}

func (m *ReaderMock) FetchResults(exerciseID int) ([]Result, error) {
	args := m.Called(exerciseID)
	return args.Get(0).([]Result), args.Error(1)
}

type WriterMock struct {
	mock.Mock
}

func (m *WriterMock) StoreResult(result *Result) error {
	args := m.Called(result)
	return args.Error(0)
}

func (m *WriterMock) UpdateExerciseProjection(exerciseID int, projection ResultsProjection) error {
	args := m.Called(exerciseID, projection)
	return args.Error(0)
}

type ProjectionBuilderMock struct {
	mock.Mock
}

func (m *ProjectionBuilderMock) Projection(results []Result) ResultsProjection {
	args := m.Called(results)
	return args.Get(0).(ResultsProjection)
}
