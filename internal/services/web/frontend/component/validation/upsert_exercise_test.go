package validation

import (
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/stretchr/testify/assert"
)

func TestValidateStoreExercise_validInsert(t *testing.T) {
	exercise := frontend.Exercise{
		Lesson: &frontend.Lesson{
			Id: 10,
		},
		Question: "question",
		Answer:   "answer",
	}

	validation := ValidateUpsertExercise(exercise, []string{"bar"})

	assert.False(t, validation.Failed())
}

func TestValidateUpsertExercise_validUpdate(t *testing.T) {
	exercise := frontend.Exercise{
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
		exercise  frontend.Exercise
		questions []string
		message   string
	}{
		{
			"empty",
			frontend.Exercise{},
			nil,
			"exercise.question is required\nexercise.answer is required\nlesson.id is required",
		},
		{
			"non unique question",
			frontend.Exercise{
				Question: "question",
			},
			[]string{"question"},
			"exercise.question must be unique\nexercise.answer is required\nlesson.id is required",
		},
		{
			"missing lesson id",
			frontend.Exercise{
				Question: "question",
				Answer:   "answer",
				Lesson:   &frontend.Lesson{},
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
