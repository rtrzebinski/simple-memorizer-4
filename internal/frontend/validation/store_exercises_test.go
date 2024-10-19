package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateStoreExercise_valid(t *testing.T) {
	exercises := models.Exercises{
		models.Exercise{
			Question: "question",
			Answer:   "answer",
		},
	}

	validation := ValidateStoreExercises(exercises)

	assert.False(t, validation.Failed())
}

func TestValidateStoreExercise_invalid(t *testing.T) {
	var tests = []struct {
		name      string
		exercises models.Exercises
		questions []string
		message   string
	}{
		{
			"missing question",
			models.Exercises{
				models.Exercise{
					Answer: "answer",
				},
			},
			nil,
			"exercise.question is required",
		},
		{
			"missing answer",
			models.Exercises{
				models.Exercise{
					Question: "question",
				},
			},
			nil,
			"exercise.answer is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := ValidateStoreExercises(tt.exercises)
			assert.Equal(t, tt.message, validator.Error())
			assert.True(t, validator.Failed())
		})
	}
}
