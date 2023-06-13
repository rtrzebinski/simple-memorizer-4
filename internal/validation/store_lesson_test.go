package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateStoreLesson_valid(t *testing.T) {
	lesson := models.Lesson{
		Name:        "name",
		Description: "description",
	}

	validator := ValidateStoreLesson(lesson, []string{"foo"})

	assert.False(t, validator.Failed())
}

func TestValidateStoreLesson_invalid(t *testing.T) {
	var tests = []struct {
		name    string
		lesson  models.Lesson
		names   []string
		message string
	}{
		{
			"empty",
			models.Lesson{},
			nil,
			"lesson.name is required",
		},
		{
			"non unique name",
			models.Lesson{
				Name: "name",
			},
			[]string{"name"},
			"lesson.name must be unique",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := ValidateStoreLesson(tt.lesson, tt.names)
			assert.Equal(t, tt.message, validator.Error())
			assert.True(t, validator.Failed())
		})
	}
}
