package memorizer

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemorizer_simpleRandomization(t *testing.T) {
	e1 := models.Exercise{Id: 1}
	e2 := models.Exercise{Id: 2}

	var exercises []models.Exercise
	exercises = append(exercises, e1)
	exercises = append(exercises, e2)

	s := Service{}
	s.Init(exercises)

	res := s.Next(models.Exercise{})

	assert.True(t, res.Question == e1.Question || res.Question == e2.Question)
}

func TestMemorizer_simpleRandomization_skipPrevious(t *testing.T) {
	e1 := models.Exercise{Id: 1}
	e2 := models.Exercise{Id: 2}

	var exercises []models.Exercise
	exercises = append(exercises, e1)
	exercises = append(exercises, e2)

	s := Service{}
	s.Init(exercises)

	res := s.Next(e2)

	assert.Equal(t, e1.Question, res.Question)
}
