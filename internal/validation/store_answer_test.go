package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateStoreAnswer_valid(t *testing.T) {
	answer := models.Answer{
		Type: models.Good,
		Exercise: &models.Exercise{
			Id: 10,
		},
	}

	validator := ValidateStoreAnswer(answer)

	assert.False(t, validator.Failed())
}

func TestValidateStoreAnswer_invalid(t *testing.T) {
	var tests = []struct {
		name    string
		answer  models.Answer
		message string
	}{
		{
			"empty",
			models.Answer{},
			"answer.type is required\nexercise.id is required",
		},
		{
			"missing exercise id",
			models.Answer{
				Exercise: &models.Exercise{},
				Type:     models.Good,
			},
			"exercise.id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := ValidateStoreAnswer(tt.answer)
			assert.Equal(t, tt.message, validator.Error())
			assert.True(t, validator.Failed())
		})
	}
}
