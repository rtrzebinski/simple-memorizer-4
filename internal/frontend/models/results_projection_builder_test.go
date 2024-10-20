package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildResultsProjection(t *testing.T) {
	yesterday := time.Now().Add(-24 * time.Hour)
	today := time.Now()

	results := Results{}

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: yesterday,
	})

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: yesterday,
	})

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: yesterday,
	})

	results = append(results, Result{
		Type:      Good,
		CreatedAt: yesterday,
	})

	projection := BuildResultsProjection(results)

	assert.Equal(t, 3, projection.BadAnswers)
	assert.Equal(t, 0, projection.BadAnswersToday)
	assert.Equal(t, yesterday, projection.LatestBadAnswer)
	assert.False(t, projection.LatestBadAnswerWasToday)
	assert.Equal(t, 1, projection.GoodAnswers)
	assert.Equal(t, 0, projection.GoodAnswersToday)
	assert.Equal(t, yesterday, projection.LatestGoodAnswer)
	assert.False(t, projection.LatestGoodAnswerWasToday)

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: today,
	})

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: today,
	})

	results = append(results, Result{
		Type:      Good,
		CreatedAt: today,
	})

	results = append(results, Result{
		Type:      Good,
		CreatedAt: today,
	})

	projection = BuildResultsProjection(results)

	assert.Equal(t, 5, projection.BadAnswers)
	assert.Equal(t, 2, projection.BadAnswersToday)
	assert.Equal(t, today, projection.LatestBadAnswer)
	assert.True(t, projection.LatestBadAnswerWasToday)
	assert.Equal(t, 3, projection.GoodAnswers)
	assert.Equal(t, 2, projection.GoodAnswersToday)
	assert.Equal(t, today, projection.LatestGoodAnswer)
	assert.True(t, projection.LatestGoodAnswerWasToday)
}
