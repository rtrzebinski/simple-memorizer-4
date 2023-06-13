package validators

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidationErr(t *testing.T) {
	err := NewValidationErr()

	assert.True(t, err.Empty())

	err.Add(errors.New("foo"))

	assert.False(t, err.Empty())
	assert.Equal(t, "foo", err.All()[0].Error())
}
