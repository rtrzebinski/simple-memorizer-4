package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateStoreResult_valid(t *testing.T) {
	result := models.Result{
		Type: models.Good,
		Exercise: &models.Exercise{
			Id: 10,
		},
	}

	validator := ValidateStoreResult(result)

	assert.False(t, validator.Failed())
}

func TestValidateStoreResult_invalid(t *testing.T) {
	var tests = []struct {
		name    string
		result  models.Result
		message string
	}{
		{
			"empty",
			models.Result{},
			"result.type is required\nexercise.id is required",
		},
		{
			"missing exercise id",
			models.Result{
				Exercise: &models.Exercise{},
				Type:     models.Good,
			},
			"exercise.id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := ValidateStoreResult(tt.result)
			assert.Equal(t, tt.message, validator.Error())
			assert.True(t, validator.Failed())
		})
	}
}
