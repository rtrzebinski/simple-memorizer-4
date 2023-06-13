package validators

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

func ValidateStoreLesson(l models.Lesson, names []string) error {
	err := NewValidationErr()

	if l.Name == "" {
		err.Add(ErrLessonNameRequired)
	} else if names != nil {
		for _, n := range names {
			if l.Name == n {
				err.Add(ErrLessonNameUnique)
			}
		}
	}

	if err.Empty() {
		return nil
	}

	return err
}
