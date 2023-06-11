package validators

import (
	"errors"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

func ValidateLessonIdentified(l models.Lesson) error {
	var err error

	if l.Id == 0 {
		err = errors.Join(err, ErrLessonIdRequired)
	}

	if err == nil {
		return nil
	}

	return NewValidationErr(err)
}
