package validation

import (
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestValidateStoreResult_valid(t *testing.T) {
	result := backend.Result{
		Type: backend.Good,
		Exercise: &backend.Exercise{
			Id: 10,
		},
	}

	validator := ValidateStoreResult(result)

	assert.False(t, validator.Failed())
}

func TestValidateStoreResult_invalid(t *testing.T) {
	var tests = []struct {
		name    string
		result  backend.Result
		message string
	}{
		{
			"empty",
			backend.Result{},
			"result.type is required\nexercise.id is required",
		},
		{
			"missing exercise id",
			backend.Result{
				Exercise: &backend.Exercise{},
				Type:     backend.Good,
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
