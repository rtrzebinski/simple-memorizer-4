package validators

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

func ValidateLessonIdentified(l models.Lesson) error {
	err := NewValidationErr()

	if l.Id == 0 {
		err.Add(ErrLessonIdRequired)
	}

	if err.Empty() {
		return nil
	}

	return err
}
