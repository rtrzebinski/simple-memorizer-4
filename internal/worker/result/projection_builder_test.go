package result

import (
	"testing"
	"time"

	"github.com/rtrzebinski/simple-memorizer-4/internal/worker"
	"github.com/stretchr/testify/assert"
)

func TestBuildProjection(t *testing.T) {
	yesterday := time.Now().Add(-24 * time.Hour)
	today := time.Now()

	var results []worker.Result

	results = append(results, worker.Result{
		Type:      worker.Bad,
		CreatedAt: yesterday,
	})

	results = append(results, worker.Result{
		Type:      worker.Bad,
		CreatedAt: yesterday,
	})

	results = append(results, worker.Result{
		Type:      worker.Bad,
		CreatedAt: yesterday,
	})

	results = append(results, worker.Result{
		Type:      worker.Good,
		CreatedAt: yesterday,
	})

	builder := NewProjectionBuilder()

	projection := builder.Projection(results)

	assert.Equal(t, 3, projection.BadAnswers)
	assert.Equal(t, 0, projection.BadAnswersToday)
	assert.Equal(t, yesterday, projection.LatestBadAnswer)
	assert.False(t, projection.LatestBadAnswerWasToday)
	assert.Equal(t, 1, projection.GoodAnswers)
	assert.Equal(t, 0, projection.GoodAnswersToday)
	assert.Equal(t, yesterday, projection.LatestGoodAnswer)
	assert.False(t, projection.LatestGoodAnswerWasToday)

	results = append(results, worker.Result{
		Type:      worker.Bad,
		CreatedAt: today,
	})

	results = append(results, worker.Result{
		Type:      worker.Bad,
		CreatedAt: today,
	})

	results = append(results, worker.Result{
		Type:      worker.Good,
		CreatedAt: today,
	})

	results = append(results, worker.Result{
		Type:      worker.Good,
		CreatedAt: today,
	})

	projection = builder.Projection(results)

	assert.Equal(t, 5, projection.BadAnswers)
	assert.Equal(t, 2, projection.BadAnswersToday)
	assert.Equal(t, today, projection.LatestBadAnswer)
	assert.True(t, projection.LatestBadAnswerWasToday)
	assert.Equal(t, 3, projection.GoodAnswers)
	assert.Equal(t, 2, projection.GoodAnswersToday)
	assert.Equal(t, today, projection.LatestGoodAnswer)
	assert.True(t, projection.LatestGoodAnswerWasToday)
}
