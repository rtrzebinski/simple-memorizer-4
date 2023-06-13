package validators

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

func ValidateStoreExercise(e models.Exercise, questions []string) error {
	err := NewValidationErr()

	if e.Question == "" {
		err.Add(ErrExerciseQuestionRequired)
	} else if questions != nil {
		for _, q := range questions {
			if e.Question == q {
				err.Add(ErrExerciseQuestionUnique)
			}
		}
	}

	if e.Answer == "" {
		err.Add(ErrExerciseAnswerRequired)
	}

	if e.Id == 0 {
		if e.Lesson == nil {
			err.Add(ErrLessonIdRequired)
		} else if e.Lesson.Id == 0 {
			err.Add(ErrLessonIdRequired)
		}
	}

	if err.Empty() {
		return nil
	}

	return err
}
