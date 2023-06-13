package validators

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

func ValidateExerciseIdentified(e models.Exercise) error {
	err := NewValidationErr()

	if e.Id == 0 {
		err.Add(ErrExerciseIdRequired)
	}

	if err.Empty() {
		return nil
	}

	return err
}
