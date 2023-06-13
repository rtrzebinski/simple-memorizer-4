package validators

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateExerciseIdentified_valid(t *testing.T) {
	exercise := models.Exercise{
		Id: 10,
	}

	err := ValidateExerciseIdentified(exercise)

	assert.NoError(t, err)
}

func TestValidateExerciseIdentified_invalid(t *testing.T) {
	exercise := models.Exercise{}

	err := ValidateExerciseIdentified(exercise)

	assert.Equal(t, "exercise.id is required", err.Error())
}
