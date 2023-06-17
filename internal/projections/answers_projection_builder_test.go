package projections

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildAnswersProjection(t *testing.T) {
	yesterday := time.Now().Add(-24 * time.Hour)
	now := time.Now()

	answers := models.Answers{}

	answers = append(answers, models.Answer{
		Type:      models.Bad,
		CreatedAt: now,
	})

	answers = append(answers, models.Answer{
		Type:      models.Bad,
		CreatedAt: yesterday,
	})

	answers = append(answers, models.Answer{
		Type:      models.Bad,
		CreatedAt: yesterday,
	})

	answers = append(answers, models.Answer{
		Type:      models.Bad,
		CreatedAt: yesterday,
	})

	answers = append(answers, models.Answer{
		Type:      models.Good,
		CreatedAt: now,
	})

	answers = append(answers, models.Answer{
		Type:      models.Good,
		CreatedAt: now,
	})

	answers = append(answers, models.Answer{
		Type:      models.Good,
		CreatedAt: yesterday,
	})

	projection := BuildAnswersProjection(answers)

	assert.Equal(t, now, projection.LatestBadAnswer)
	assert.Equal(t, 4, projection.BadAnswers)
	assert.Equal(t, 1, projection.BadAnswersToday)
	assert.Equal(t, now, projection.LatestGoodAnswer)
	assert.Equal(t, 3, projection.GoodAnswers)
	assert.Equal(t, 2, projection.GoodAnswersToday)
}
