package worker

import (
	"context"
	"time"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
)

func (s *PostgresSuite) TestReader_FetchResults() {
	ctx := context.Background()

	db := s.db

	r := NewReader(db)

	exercise := &postgres.Exercise{}
	postgres.CreateExercise(db, exercise)

	postgres.CreateResult(db, &postgres.Result{
		Type:       "good",
		ExerciseId: exercise.Id,
	})

	results, err := r.FetchResults(ctx, exercise.Id)
	s.NoError(err)
	s.Len(results, 1)
	s.Equal(1, results[0].Id)
	s.Equal(worker.ResultType("good"), results[0].Type)
	s.Equal(exercise.Id, results[0].ExerciseId)
	s.Equal(time.Now().Local().Format("2006-01-02"), results[0].CreatedAt.Local().Format("2006-01-02"))
}
