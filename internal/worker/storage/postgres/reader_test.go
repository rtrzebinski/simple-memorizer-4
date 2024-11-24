package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/worker"
	"time"
)

func (suite *PostgresSuite) TestReader_FetchResults() {
	ctx := context.Background()

	db := suite.db

	r := NewReader(db)

	exercise := &Exercise{}
	createExercise(db, exercise)

	createResult(db, &Result{
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
