package worker

import (
	"context"
	"time"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
)

func (suite *PostgresSuite) TestReader_FetchResults() {
	ctx := context.Background()

	db := suite.db

	r := NewReader(db)

	exercise := &postgres.Exercise{}
	postgres.CreateExercise(db, exercise)

	postgres.CreateResult(db, &postgres.Result{
		Type:       "good",
		ExerciseId: exercise.Id,
	})

	results, err := r.FetchResults(ctx, exercise.Id)
	suite.NoError(err)
	suite.Len(results, 1)
	suite.Equal(1, results[0].Id)
	suite.Equal(worker.ResultType("good"), results[0].Type)
	suite.Equal(exercise.Id, results[0].ExerciseId)
	suite.Equal(time.Now().Local().Format("2006-01-02"), results[0].CreatedAt.Local().Format("2006-01-02"))
}
