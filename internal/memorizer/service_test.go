package memorizer

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
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
		Id: 2,
		AnswersProjection: models.AnswersProjection{
			GoodAnswers: 10,
		},
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
	assert.Equal(t, e2.AnswersProjection.GoodAnswers, res.AnswersProjection.GoodAnswers)
}

func TestMemorizer_points(t *testing.T) {
	var tests = []struct {
		percent int
		points  int
	}{
		{0, 10},
		{5, 10},
		{10, 10},
		{15, 9},
		{20, 9},
		{25, 8},
		{30, 8},
		{35, 7},
		{40, 7},
		{45, 6},
		{50, 6},
		{55, 5},
		{60, 5},
		{65, 4},
		{70, 4},
		{75, 3},
		{80, 3},
		{85, 2},
		{90, 2},
		{95, 1},
		{100, 1},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.points, points(tt.percent))
		})
	}
}
