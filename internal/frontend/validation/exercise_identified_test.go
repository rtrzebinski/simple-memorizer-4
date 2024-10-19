package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateExerciseIdentified_valid(t *testing.T) {
	exercise := models.Exercise{
		Id: 10,
	}

	validator := ValidateExerciseIdentified(exercise)

	assert.False(t, validator.Failed())
}

func TestValidateExerciseIdentified_invalid(t *testing.T) {
	exercise := models.Exercise{}

	validator := ValidateExerciseIdentified(exercise)

	assert.Equal(t, "exercise.id is required", validator.Error())
	assert.True(t, validator.Failed())
}
