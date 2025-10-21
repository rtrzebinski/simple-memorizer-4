package postgres

import (
	"testing"
	"time"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
	"github.com/stretchr/testify/suite"
)

type WorkerReaderSuite struct {
	PostgresSuite
	reader *WorkerReader
}

func TestWorkerReader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(WorkerReaderSuite))
}

func (s *WorkerReaderSuite) SetupSuite() {
	s.PostgresSuite.SetupSuite()
	s.reader = NewWorkerReader(s.DB)
}

func (s *WorkerReaderSuite) TestWorkerReader_FetchResults() {
	ctx := s.T().Context()

	exercise := &exercise{}
	createExercise(s.DB, exercise)

	createResult(s.DB, &result{
		Type:       "good",
		ExerciseId: exercise.Id,
	})

	results, err := s.reader.FetchResults(ctx, exercise.Id)
	s.NoError(err)
	s.Len(results, 1)
	s.Equal(1, results[0].Id)
	s.Equal(worker.ResultType("good"), results[0].Type)
	s.Equal(exercise.Id, results[0].ExerciseId)
	s.Equal(time.Now().Local().Format("2006-01-02"), results[0].CreatedAt.Local().Format("2006-01-02"))
}
