package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

func ValidateLessonIdentified(l models.Lesson) *Validator {
	validator := NewValidator()

	if l.Id == 0 {
		validator.AddError(ErrLessonIdRequired)
	}

	return validator
}
