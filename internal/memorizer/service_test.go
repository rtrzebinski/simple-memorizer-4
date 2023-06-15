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
	e2 := models.Exercise{Id: 2, GoodAnswers: 10}

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
