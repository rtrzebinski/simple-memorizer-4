package validators

import (
	"errors"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

func ValidateStoreExercise(e models.Exercise, questions []string) error {
	var err error

	if e.Question == "" {
		err = errors.Join(err, ErrExerciseQuestionRequired)
	} else if questions != nil {
		for _, q := range questions {
			if e.Question == q {
				err = errors.Join(err, ErrExerciseQuestionUnique)
			}
		}
	}

	if e.Answer == "" {
		err = errors.Join(err, ErrExerciseAnswerRequired)
	}

	if e.Id == 0 {
		if e.Lesson == nil {
			err = errors.Join(err, ErrLessonIdRequired)
		} else if e.Lesson.Id == 0 {
			err = errors.Join(err, ErrLessonIdRequired)
		}
	}

	if err == nil {
		return nil
	}

	return NewValidationErr(err)
}
