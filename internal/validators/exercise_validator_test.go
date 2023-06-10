package validators

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateStoreExercise_validInsert(t *testing.T) {
	exercise := models.Exercise{
		Lesson: &models.Lesson{
			Id: 10,
		},
		Question: "question",
		Answer:   "answer",
	}

	err := ValidateStoreExercise(exercise, []string{"bar"})

	assert.NoError(t, err)
}

func TestValidateStoreExercise_validUpdate(t *testing.T) {
	exercise := models.Exercise{
		Id:       10,
		Question: "question",
		Answer:   "answer",
	}

	err := ValidateStoreExercise(exercise, []string{"bar"})

	assert.NoError(t, err)
}

func TestValidateStoreExercise_invalid(t *testing.T) {
	var tests = []struct {
		name      string
		exercise  models.Exercise
		questions []string
	}{
		{
			"empty",
			models.Exercise{},
			nil,
		},
		{
			"non unique question",
			models.Exercise{
				Question: "question",
			},
			[]string{"question"},
		},
		{
			"missing lesson id",
			models.Exercise{
				Question: "question",
				Answer:   "answer",
				Lesson:   &models.Lesson{},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := ValidateStoreExercise(tt.exercise, tt.questions)
			assert.True(t, IsValidationErr(err))
		})
	}
}
