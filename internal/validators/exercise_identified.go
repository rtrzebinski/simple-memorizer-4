package validators

import (
	"errors"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

func ValidateExerciseIdentified(e models.Exercise) error {
	var err error

	if e.Id == 0 {
		err = errors.Join(err, ErrExerciseIdRequired)
	}

	if err == nil {
		return nil
	}

	return NewValidationErr(err)
}
