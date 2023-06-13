package validators

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateLessonIdentified_valid(t *testing.T) {
	lesson := models.Lesson{
		Id: 10,
	}

	err := ValidateLessonIdentified(lesson)

	assert.NoError(t, err)
}

func TestValidateLessonIdentified_invalid(t *testing.T) {
	lesson := models.Lesson{}

	err := ValidateLessonIdentified(lesson)

	assert.Equal(t, "lesson.id is required", err.Error())
}
