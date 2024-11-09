package validation

import (
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestValidateUpsertLesson_valid(t *testing.T) {
	lesson := backend.Lesson{
		Name:        "name",
		Description: "description",
	}

	validator := ValidateUpsertLesson(lesson, []string{"foo"})

	assert.False(t, validator.Failed())
}

func TestValidateUpsertLesson_invalid(t *testing.T) {
	var tests = []struct {
		name    string
		lesson  backend.Lesson
		names   []string
		message string
	}{
		{
			"empty",
			backend.Lesson{},
			nil,
			"lesson.name is required",
		},
		{
			"non unique name",
			backend.Lesson{
				Name: "name",
			},
			[]string{"name"},
			"lesson.name must be unique",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := ValidateUpsertLesson(tt.lesson, tt.names)
			assert.Equal(t, tt.message, validator.Error())
			assert.True(t, validator.Failed())
		})
	}
}
