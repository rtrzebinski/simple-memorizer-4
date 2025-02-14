package validation

import "github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"

func ValidateExerciseIdentified(e backend.Exercise) *Validator {
	validator := NewValidator()

	if e.Id == 0 {
		validator.AddError(ErrExerciseIdRequired)
	}

	return validator
}
