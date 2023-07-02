package projections

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildResultsProjection(t *testing.T) {
	yesterday := time.Now().Add(-24 * time.Hour)
	today := time.Now()

	results := models.Results{}

	results = append(results, models.Result{
		Type:      models.Bad,
		CreatedAt: yesterday,
	})

	results = append(results, models.Result{
		Type:      models.Bad,
		CreatedAt: yesterday,
	})

	results = append(results, models.Result{
		Type:      models.Bad,
		CreatedAt: yesterday,
	})

	results = append(results, models.Result{
		Type:      models.Good,
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

	results = append(results, models.Result{
		Type:      models.Bad,
		CreatedAt: today,
	})

	results = append(results, models.Result{
		Type:      models.Bad,
		CreatedAt: today,
	})

	results = append(results, models.Result{
		Type:      models.Good,
		CreatedAt: today,
	})

	results = append(results, models.Result{
		Type:      models.Good,
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
