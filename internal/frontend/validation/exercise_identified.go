package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
)

func ValidateExerciseIdentified(e models.Exercise) *Validator {
	validator := NewValidator()

	if e.Id == 0 {
		validator.AddError(ErrExerciseIdRequired)
	}

	return validator
}
