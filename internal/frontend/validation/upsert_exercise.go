package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
)

func ValidateUpsertExercise(e models.Exercise, questions []string) *Validator {
	validator := NewValidator()

	if e.Question == "" {
		validator.AddError(ErrExerciseQuestionRequired)
	} else if questions != nil {
		for _, q := range questions {
			if e.Question == q {
				validator.AddError(ErrExerciseQuestionUnique)
			}
		}
	}

	if e.Answer == "" {
		validator.AddError(ErrExerciseAnswerRequired)
	}

	if e.Id == 0 {
		if e.Lesson == nil {
			validator.AddError(ErrLessonIdRequired)
		} else if e.Lesson.Id == 0 {
			validator.AddError(ErrLessonIdRequired)
		}
	}

	return validator
}
