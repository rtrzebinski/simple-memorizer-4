package validation

import (
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"github.com/stretchr/testify/assert"
)

func TestValidateUpsertLesson_valid(t *testing.T) {
	lesson := frontend.Lesson{
		Name:        "name",
		Description: "description",
	}

	validator := ValidateUpsertLesson(lesson, []string{"foo"})

	assert.False(t, validator.Failed())
}

func TestValidateUpsertLesson_invalid(t *testing.T) {
	var tests = []struct {
		name    string
		lesson  frontend.Lesson
		names   []string
		message string
	}{
		{
			"empty",
			frontend.Lesson{},
			nil,
			"lesson.name is required",
		},
		{
			"non unique name",
			frontend.Lesson{
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
