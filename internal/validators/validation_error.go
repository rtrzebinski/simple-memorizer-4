package validators

import "errors"

// ValidationErr error type for validation.
type ValidationErr struct{ err error }

// NewValidationErr creates a new validation error.
func NewValidationErr(err error) error {
	return ValidationErr{err: err}
}

// Error implements the error interface function.
func (v ValidationErr) Error() string { return v.err.Error() }

// IsValidationErr checks if the provided error is from the type `ValidationErr`.
func IsValidationErr(err error) bool {
	var validationErr ValidationErr

	return errors.As(err, &validationErr)
}
