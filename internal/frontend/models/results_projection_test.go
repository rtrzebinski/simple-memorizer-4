package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GoodAnswersPercent(t *testing.T) {
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
			projection := ResultsProjection{BadAnswers: tt.bad, GoodAnswers: tt.good}
			assert.Equal(t, tt.percent, projection.GoodAnswersPercent())
		})
	}
}

func Test_RegisterGoodAnswer(t *testing.T) {
	projection := ResultsProjection{}

	assert.Equal(t, 0, projection.GoodAnswers)
	assert.Equal(t, 0, projection.GoodAnswersToday)
	assert.Equal(t, false, projection.LatestGoodAnswerWasToday)
	assert.Empty(t, projection.LatestGoodAnswer)

	projection.RegisterGoodAnswer()

	assert.Equal(t, 1, projection.GoodAnswers)
	assert.Equal(t, 1, projection.GoodAnswersToday)
	assert.Equal(t, true, projection.LatestGoodAnswerWasToday)
	assert.Equal(t, time.Now().Day(), projection.LatestGoodAnswer.Day())
}

func Test_RegisterBadAnswer(t *testing.T) {
	projection := ResultsProjection{}

	assert.Equal(t, 0, projection.BadAnswers)
	assert.Equal(t, 0, projection.BadAnswersToday)
	assert.Equal(t, false, projection.LatestBadAnswerWasToday)
	assert.Empty(t, projection.LatestBadAnswer)

	projection.RegisterBadAnswer()

	assert.Equal(t, 1, projection.BadAnswers)
	assert.Equal(t, 1, projection.BadAnswersToday)
	assert.Equal(t, true, projection.LatestBadAnswerWasToday)
	assert.Equal(t, time.Now().Day(), projection.LatestBadAnswer.Day())
}
