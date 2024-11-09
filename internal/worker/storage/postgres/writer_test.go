package postgres

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/worker"
	"github.com/stretchr/testify/assert"
)

func (suite *PostgresSuite) TestWriter_StoreAnswer() {
	db := suite.db

	w := NewWriter(db)

	exercise := &Exercise{}
	createExercise(db, exercise)

	result := worker.Result{
		Type:       worker.Good,
		ExerciseId: exercise.Id,
	}

	err := w.StoreResult(&result)
	assert.NoError(suite.T(), err)

	stored := fetchLatestResult(db)

	assert.Equal(suite.T(), string(result.Type), stored.Type)
	assert.Equal(suite.T(), result.ExerciseId, stored.ExerciseId)
}
