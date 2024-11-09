package validation

import (
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestValidateExerciseIdentified_valid(t *testing.T) {
	exercise := backend.Exercise{
		Id: 10,
	}

	validator := ValidateExerciseIdentified(exercise)

	assert.False(t, validator.Failed())
}

func TestValidateExerciseIdentified_invalid(t *testing.T) {
	exercise := backend.Exercise{}

	validator := ValidateExerciseIdentified(exercise)

	assert.Equal(t, "exercise.id is required", validator.Error())
	assert.True(t, validator.Failed())
}
