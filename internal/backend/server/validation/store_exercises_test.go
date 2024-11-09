package validation

import (
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestValidateStoreExercise_valid(t *testing.T) {
	exercises := backend.Exercises{
		backend.Exercise{
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
		exercises backend.Exercises
		questions []string
		message   string
	}{
		{
			"missing question",
			backend.Exercises{
				backend.Exercise{
					Answer: "answer",
				},
			},
			nil,
			"exercise.question is required",
		},
		{
			"missing answer",
			backend.Exercises{
				backend.Exercise{
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
