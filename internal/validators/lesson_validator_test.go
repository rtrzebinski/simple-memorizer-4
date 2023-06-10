package validators

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

	err := ValidateStoreLesson(lesson, []string{"foo"})

	assert.NoError(t, err)
}

func TestValidateStoreLesson_invalid(t *testing.T) {
	var tests = []struct {
		name   string
		lesson models.Lesson
		names  []string
	}{
		{
			"empty",
			models.Lesson{},
			nil,
		},
		{
			"non unique name",
			models.Lesson{
				Name: "name",
			},
			[]string{"name"},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := ValidateStoreLesson(tt.lesson, tt.names)
			assert.True(t, IsValidationErr(err))
		})
	}
}
