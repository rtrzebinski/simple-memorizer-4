package validators

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidationErr(t *testing.T) {
	err := NewValidationErr(errors.New("foo"))

	assert.True(t, IsValidationErr(err))
	assert.Equal(t, "foo", err.Error())
}
