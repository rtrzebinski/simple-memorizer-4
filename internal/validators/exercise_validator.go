package validators

import (
	"errors"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

var (
	ErrExerciseQuestionRequired = errors.New("exercise.question is required")
	ErrExerciseQuestionUnique   = errors.New("exercise.question must be unique")
	ErrExerciseAnswerRequired   = errors.New("exercise.answer is required")
	ErrLessonIdRequired         = errors.New("lesson.id is required")
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
