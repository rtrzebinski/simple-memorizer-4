package validators

import (
	"errors"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
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
