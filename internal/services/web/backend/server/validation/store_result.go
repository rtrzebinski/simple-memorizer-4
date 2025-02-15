package validation

import "github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"

func ValidateStoreResult(a backend.Result) *Validator {
	validator := NewValidator()

	if a.Type == "" {
		validator.AddError(ErrResultTypeRequired)
	}

	if a.Exercise == nil {
		validator.AddError(ErrExerciseIdRequired)
	} else if a.Exercise.Id == 0 {
		validator.AddError(ErrExerciseIdRequired)
	}

	return validator
}
