package validation

import "github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"

func ValidateStoreExercises(exercises backend.Exercises) *Validator {
	validator := NewValidator()

	for _, e := range exercises {
		if e.Question == "" {
			validator.AddError(ErrExerciseQuestionRequired)
		}

		if e.Answer == "" {
			validator.AddError(ErrExerciseAnswerRequired)
		}
	}

	return validator
}
