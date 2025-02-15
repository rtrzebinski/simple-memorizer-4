package validation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	validation := NewValidator()

	assert.False(t, validation.Failed())

	validation.AddError(errors.New("foo"))

	assert.True(t, validation.Failed())
	assert.Equal(t, "foo", validation.Errors()[0].Error())
	assert.Equal(t, "foo", validation.Error())

	validation.AddError(errors.New("bar"))
	assert.Equal(t, "bar", validation.Errors()[1].Error())
	assert.Equal(t, "foo\nbar", validation.Error())
}
