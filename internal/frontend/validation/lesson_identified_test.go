package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateLessonIdentified_valid(t *testing.T) {
	lesson := models.Lesson{
		Id: 10,
	}

	validator := ValidateLessonIdentified(lesson)

	assert.False(t, validator.Failed())
}

func TestValidateLessonIdentified_invalid(t *testing.T) {
	lesson := models.Lesson{}

	validator := ValidateLessonIdentified(lesson)

	assert.Equal(t, "lesson.id is required", validator.Error())
	assert.True(t, validator.Failed())
}
