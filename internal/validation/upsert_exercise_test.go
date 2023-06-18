package validation

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateStoreExercise_validInsert(t *testing.T) {
	exercise := models.Exercise{
		Lesson: &models.Lesson{
			Id: 10,
		},
		Question: "question",
		Answer:   "answer",
	}

	validation := ValidateUpsertExercise(exercise, []string{"bar"})

	assert.False(t, validation.Failed())
}

func TestValidateUpsertExercise_validUpdate(t *testing.T) {
	exercise := models.Exercise{
		Id:       10,
		Question: "question",
		Answer:   "answer",
	}

	validator := ValidateUpsertExercise(exercise, []string{"bar"})

	assert.False(t, validator.Failed())
}

func TestValidateUpsertExercise_invalid(t *testing.T) {
	var tests = []struct {
		name      string
		exercise  models.Exercise
		questions []string
		message   string
	}{
		{
			"empty",
			models.Exercise{},
			nil,
			"exercise.question is required\nexercise.answer is required\nlesson.id is required",
		},
		{
			"non unique question",
			models.Exercise{
				Question: "question",
			},
			[]string{"question"},
			"exercise.question must be unique\nexercise.answer is required\nlesson.id is required",
		},
		{
			"missing lesson id",
			models.Exercise{
				Question: "question",
				Answer:   "answer",
				Lesson:   &models.Lesson{},
			},
			nil,
			"lesson.id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := ValidateUpsertExercise(tt.exercise, tt.questions)
			assert.Equal(t, tt.message, validator.Error())
			assert.True(t, validator.Failed())
		})
	}
}
