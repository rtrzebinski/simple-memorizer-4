package frontend

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/projections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	service *Service
	reader  *ReaderMock
	writer  *WriterMock
}

func (suite *ServiceSuite) SetupTest() {
	suite.reader = NewReaderMock()
	suite.writer = NewWriterMock()
	suite.service = NewService(suite.reader, suite.writer)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (suite *ServiceSuite) TestFetchLessons() {
	expectedLessons := models.Lessons{{Id: 1, Name: "Lesson 1"}, {Id: 2, Name: "Lesson 2"}}
	suite.reader.On("FetchLessons").Return(expectedLessons, nil)

	lessons, err := suite.service.FetchLessons()

	suite.reader.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), lessons, len(expectedLessons))
}

func (suite *ServiceSuite) TestHydrateLesson() {
	lesson := &models.Lesson{Id: 1, Name: "Lesson 1"}
	suite.reader.On("HydrateLesson", lesson).Return(nil)

	err := suite.service.HydrateLesson(lesson)

	suite.reader.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
}

func (suite *ServiceSuite) TestFetchExercises() {
	lesson := models.Lesson{Id: 1}
	exercises := models.Exercises{
		{Id: 1, Results: []models.Result{{Type: models.Good}, {Type: models.Good}}},
		{Id: 2, Results: []models.Result{{Type: models.Bad}, {Type: models.Bad}}},
	}

	expectedExercises := models.Exercises{
		{Id: 1, Results: []models.Result{{Type: models.Good}, {Type: models.Good}},
			ResultsProjection: projections.BuildResultsProjection(exercises[0].Results)},
		{Id: 2, Results: []models.Result{{Type: models.Bad}, {Type: models.Bad}},
			ResultsProjection: projections.BuildResultsProjection(exercises[1].Results)},
	}

	suite.reader.On("FetchExercises", lesson).Return(exercises, nil)

	result, err := suite.service.FetchExercises(lesson)

	suite.Nil(err)
	suite.Equal(expectedExercises, result)
	suite.Equal(expectedExercises[0].ResultsProjection, result[0].ResultsProjection)
	suite.Equal(expectedExercises[1].ResultsProjection, result[1].ResultsProjection)
	suite.reader.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestUpsertLesson() {
	lesson := &models.Lesson{Id: 1, Name: "Lesson 1"}

	suite.writer.On("UpsertLesson", lesson).Return(nil)

	err := suite.service.UpsertLesson(lesson)

	suite.writer.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
}

func (suite *ServiceSuite) TestDeleteLesson() {
	lesson := models.Lesson{Id: 1, Name: "Lesson 1"}

	suite.writer.On("DeleteLesson", lesson).Return(nil)

	err := suite.service.DeleteLesson(lesson)

	suite.writer.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
}

func (suite *ServiceSuite) TestUpsertExercise() {
	exercise := &models.Exercise{Id: 1, Question: "Exercise 1"}

	suite.writer.On("UpsertExercise", exercise).Return(nil)

	err := suite.service.UpsertExercise(exercise)

	suite.writer.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
}

func (suite *ServiceSuite) TestStoreExercises() {
	exercises := models.Exercises{
		models.Exercise{Id: 1, Question: "Exercise 1"},
	}

	suite.writer.On("StoreExercises", exercises).Return(nil)

	err := suite.service.StoreExercises(exercises)

	suite.writer.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
}

func (suite *ServiceSuite) TestDeleteExercise() {
	exercise := models.Exercise{Id: 1, Question: "Exercise 1"}

	suite.writer.On("DeleteExercise", exercise).Return(nil)

	err := suite.service.DeleteExercise(exercise)

	suite.writer.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
}

func (suite *ServiceSuite) TestStoreResult() {
	result := &models.Result{Id: 1, Type: models.Good}

	suite.writer.On("StoreResult", result).Return(nil)

	err := suite.service.StoreResult(result)

	suite.writer.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
}
