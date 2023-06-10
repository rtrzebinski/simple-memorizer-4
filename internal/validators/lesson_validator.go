package validators

import (
	"errors"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

var (
	ErrLessonNameRequired = errors.New("lesson.name is required")
	ErrLessonNameUnique   = errors.New("lesson.name must be unique")
)

func ValidateStoreLesson(l models.Lesson, names []string) error {
	var err error

	if l.Name == "" {
		err = errors.Join(err, ErrLessonNameRequired)
	} else if names != nil {
		for _, n := range names {
			if l.Name == n {
				err = errors.Join(err, ErrLessonNameUnique)
			}
		}
	}

	if err == nil {
		return nil
	}

	return NewValidationErr(err)
}
