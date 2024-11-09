package validation

import "github.com/rtrzebinski/simple-memorizer-4/internal/backend"

func ValidateUpsertLesson(l backend.Lesson, names []string) *Validator {
	validator := NewValidator()

	if l.Name == "" {
		validator.AddError(ErrLessonNameRequired)
	} else if names != nil {
		for _, n := range names {
			if l.Name == n {
				validator.AddError(ErrLessonNameUnique)
			}
		}
	}

	return validator
}
