package worker

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ReaderMock struct {
	mock.Mock
}

func (m *ReaderMock) FetchResults(ctx context.Context, exerciseID int) ([]Result, error) {
	args := m.Called(ctx, exerciseID)
	return args.Get(0).([]Result), args.Error(1)
}

type WriterMock struct {
	mock.Mock
}

func (m *WriterMock) StoreResult(ctx context.Context, result Result) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *WriterMock) UpdateExerciseProjection(ctx context.Context, exerciseID int, projection ResultsProjection) error {
	args := m.Called(ctx, exerciseID, projection)
	return args.Error(0)
}
