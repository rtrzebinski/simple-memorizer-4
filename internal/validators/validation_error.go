package validators

import (
	"errors"
)

// ValidationErr error type for validation.
type ValidationErr struct {
	errors []error
}

// NewValidationErr creates a new validation error.
func NewValidationErr() *ValidationErr {
	return &ValidationErr{}
}

// Add adds a new error.
func (v *ValidationErr) Add(err error) {
	v.errors = append(v.errors, err)
}

// All retrieves all errors.
func (v *ValidationErr) All() []error {
	return v.errors
}

// Empty checks if there are no errors.
func (v *ValidationErr) Empty() bool {
	return v.errors == nil
}

// Error implements the error interface function.
func (v *ValidationErr) Error() string {
	var out error

	for _, e := range v.errors {
		out = errors.Join(out, e)
	}

	return out.Error()
}
