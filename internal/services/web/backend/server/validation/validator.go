package validation

import (
	"errors"
)

type Validator struct {
	errors []error
}

// NewValidator creates a new validation error.
func NewValidator() *Validator {
	return &Validator{}
}

// AddError adds a new error.
func (v *Validator) AddError(err error) {
	v.errors = append(v.errors, err)
}

// Failed checks if there are any errors.
func (v *Validator) Failed() bool {
	return v.errors != nil
}

// Errors retrieves all errors.
func (v *Validator) Errors() []error {
	return v.errors
}

// Error implements the error interface function.
func (v *Validator) Error() string {
	var out error

	for _, e := range v.errors {
		out = errors.Join(out, e)
	}

	return out.Error()
}
