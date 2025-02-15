package frontend

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExercise_GoodAnswersPercent(t *testing.T) {
	var tests = []struct {
		good    int
		bad     int
		percent int
	}{
		{0, 0, 0},
		{1, 19, 5},
		{1, 11, 8},
		{1, 9, 10},
		{1, 1, 50},
		{19, 1, 95},
		{10, 0, 100},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			e := Exercise{BadAnswers: tt.bad, GoodAnswers: tt.good}
			assert.Equal(t, tt.percent, e.GoodAnswersPercent())
		})
	}
}

func TestExercise_RegisterGoodAnswer(t *testing.T) {
	exercise := Exercise{}

	assert.Equal(t, 0, exercise.GoodAnswers)
	assert.Equal(t, 0, exercise.GoodAnswersToday)
	assert.Equal(t, false, exercise.LatestGoodAnswerWasToday)
	assert.Empty(t, exercise.LatestGoodAnswer)

	exercise.RegisterGoodAnswer()

	assert.Equal(t, 1, exercise.GoodAnswers)
	assert.Equal(t, 1, exercise.GoodAnswersToday)
	assert.Equal(t, true, exercise.LatestGoodAnswerWasToday)
	assert.Equal(t, time.Now().Day(), exercise.LatestGoodAnswer.Time.Day())
}

func TestExercise_RegisterBadAnswer(t *testing.T) {
	exercise := Exercise{}

	assert.Equal(t, 0, exercise.BadAnswers)
	assert.Equal(t, 0, exercise.BadAnswersToday)
	assert.Equal(t, false, exercise.LatestBadAnswerWasToday)
	assert.Empty(t, exercise.LatestBadAnswer)

	exercise.RegisterBadAnswer()

	assert.Equal(t, 1, exercise.BadAnswers)
	assert.Equal(t, 1, exercise.BadAnswersToday)
	assert.Equal(t, true, exercise.LatestBadAnswerWasToday)
	assert.Equal(t, time.Now().Day(), exercise.LatestBadAnswer.Time.Day())
}
