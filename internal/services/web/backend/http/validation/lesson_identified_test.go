package validation

import (
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/stretchr/testify/assert"
)

func TestValidateLessonIdentified_valid(t *testing.T) {
	lesson := backend.Lesson{
		Id: 10,
	}

	validator := ValidateLessonIdentified(lesson)

	assert.False(t, validator.Failed())
}

func TestValidateLessonIdentified_invalid(t *testing.T) {
	lesson := backend.Lesson{}

	validator := ValidateLessonIdentified(lesson)

	assert.Equal(t, "lesson.id is required", validator.Error())
	assert.True(t, validator.Failed())
}
