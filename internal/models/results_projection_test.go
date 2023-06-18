package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodAnswersPercent(t *testing.T) {
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
