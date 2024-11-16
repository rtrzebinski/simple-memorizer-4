package memorizer

import (
	"github.com/guregu/null/v5"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemorizer_simpleRandomization(t *testing.T) {
	e1 := models.Exercise{Id: 1}
	e2 := models.Exercise{Id: 2}

	var exercises = make(map[int]models.Exercise)
	exercises[e1.Id] = e1
	exercises[e2.Id] = e2

	s := Service{}
	s.Init(exercises)

	res := s.Next(models.Exercise{})

	assert.True(t, res.Question == e1.Question || res.Question == e2.Question)
}

func TestMemorizer_simpleRandomization_skipPrevious(t *testing.T) {
	e1 := models.Exercise{Id: 1}
	e2 := models.Exercise{
		Id:          2,
		GoodAnswers: 10,
	}

	var exercises = make(map[int]models.Exercise)
	exercises[e1.Id] = e1
	exercises[e2.Id] = e2

	s := Service{}
	s.Init(exercises)

	res := s.Next(e2)
	assert.Equal(t, e1.Question, res.Question)

	res = s.Next(e1)
	assert.Equal(t, e2.Question, res.Question)
	assert.Equal(t, e2.GoodAnswers, res.GoodAnswers)
}

func TestPoints(t *testing.T) {
	var tests = []struct {
		name           string
		exercise       models.Exercise
		expectedResult int
	}{
		{
			name: "Only good answers today",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: true,
				LatestBadAnswerWasToday:  false,
			},
			expectedResult: 1,
		},
		{
			name: "Good and bad answers today, good answer most recent",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: true,
				LatestBadAnswerWasToday:  true,
				LatestGoodAnswer:         null.TimeFrom(time.Now()),
				LatestBadAnswer:          null.TimeFrom(time.Now().Add(-time.Hour)),
			},
			expectedResult: 1,
		},
		{
			name: "Bad answers today (BadAnswersToday = 1)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  true,
				BadAnswersToday:          1,
			},
			expectedResult: 80,
		},
		{
			name: "Bad answers today (BadAnswersToday = 2)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  true,
				BadAnswersToday:          2,
			},
			expectedResult: 60,
		},
		{
			name: "Bad answers today (BadAnswersToday = 3)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  true,
				BadAnswersToday:          3,
			},
			expectedResult: 40,
		},
		{
			name: "Bad answers today (BadAnswersToday = 4)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  true,
				BadAnswersToday:          4,
			},
			expectedResult: 20,
		},
		{
			name: "Bad answers today (BadAnswersToday = 5)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  true,
				BadAnswersToday:          5,
			},
			expectedResult: 1,
		},
		{
			name: "Only good answers before today (GoodAnswers = 1)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  false,
				GoodAnswers:              1,
			},
			expectedResult: 80,
		},
		{
			name: "Only good answers before today (GoodAnswers = 2)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  false,
				GoodAnswers:              2,
			},
			expectedResult: 60,
		},
		{
			name: "Only good answers before today (GoodAnswers = 3)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  false,
				GoodAnswers:              3,
			},
			expectedResult: 40,
		},
		{
			name: "Only good answers before today (GoodAnswers = 4)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  false,
				GoodAnswers:              4,
			},
			expectedResult: 20,
		},
		{
			name: "Only good answers before today (GoodAnswers = 5)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  false,
				GoodAnswers:              5,
			},
			expectedResult: 1,
		},
		{
			name: "Only good answers before today (GoodAnswers = 1)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				BadAnswers:               0,
				GoodAnswers:              1,
			},
			expectedResult: 80,
		},
		{
			name: "Other cases (GoodAnswersPercent = 66)",
			exercise: models.Exercise{
				LatestGoodAnswerWasToday: false,
				LatestBadAnswerWasToday:  false,
				BadAnswers:               1,
				GoodAnswers:              2,
			},
			expectedResult: 34,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedResult, points(tt.exercise))
		})
	}
}
