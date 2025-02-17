package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

func ValidateLessonIdentified(l backend.Lesson) *Validator {
	validator := NewValidator()

	if l.Id == 0 {
		validator.AddError(ErrLessonIdRequired)
	}

	return validator
}
