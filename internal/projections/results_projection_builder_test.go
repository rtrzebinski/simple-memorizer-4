package projections

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildResultsProjection(t *testing.T) {
	yesterday := time.Now().Add(-24 * time.Hour)
	now := time.Now()

	results := models.Results{}

	results = append(results, models.Result{
		Type:      models.Bad,
		CreatedAt: now,
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
		Type:      models.Bad,
		CreatedAt: yesterday,
	})

	results = append(results, models.Result{
		Type:      models.Good,
		CreatedAt: now,
	})

	results = append(results, models.Result{
		Type:      models.Good,
		CreatedAt: now,
	})

	results = append(results, models.Result{
		Type:      models.Good,
		CreatedAt: yesterday,
	})

	projection := BuildResultsProjection(results)

	assert.Equal(t, now, projection.LatestBadAnswer)
	assert.Equal(t, 4, projection.BadAnswers)
	assert.Equal(t, 1, projection.BadAnswersToday)
	assert.Equal(t, now, projection.LatestGoodAnswer)
	assert.Equal(t, 3, projection.GoodAnswers)
	assert.Equal(t, 2, projection.GoodAnswersToday)
}
