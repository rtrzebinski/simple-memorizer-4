package internal

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/projections"
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

func (suite *ServiceSuite) TestFetchAllLessons() {
	expectedLessons := models.Lessons{{Id: 1, Name: "Lesson 1"}, {Id: 2, Name: "Lesson 2"}}
	suite.reader.On("FetchAllLessons").Return(expectedLessons, nil)

	lessons, err := suite.service.FetchAllLessons()

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

func (suite *ServiceSuite) TestFetchExercisesOfLesson() {
	lesson := models.Lesson{Id: 1}
	exercises := models.Exercises{
		{Id: 1, Answers: []models.Answer{{Type: models.Good}, {Type: models.Good}}},
		{Id: 2, Answers: []models.Answer{{Type: models.Bad}, {Type: models.Bad}}},
	}

	expectedExercises := models.Exercises{
		{Id: 1, Answers: []models.Answer{{Type: models.Good}, {Type: models.Good}},
			AnswersProjection: projections.BuildAnswersProjection(exercises[0].Answers)},
		{Id: 2, Answers: []models.Answer{{Type: models.Bad}, {Type: models.Bad}},
			AnswersProjection: projections.BuildAnswersProjection(exercises[1].Answers)},
	}

	suite.reader.On("FetchExercisesOfLesson", lesson).Return(exercises, nil)

	result, err := suite.service.FetchExercisesOfLesson(lesson)

	suite.Nil(err)
	suite.Equal(expectedExercises, result)
	suite.Equal(expectedExercises[0].AnswersProjection, result[0].AnswersProjection)
	suite.Equal(expectedExercises[1].AnswersProjection, result[1].AnswersProjection)
	suite.reader.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestFetchAnswersOfExercise() {
	exercise := models.Exercise{Id: 1, Question: "Exercise 1"}

	answers := models.Answers{
		{Id: 1, Type: models.Good},
		{Id: 2, Type: models.Bad},
	}

	suite.reader.On("FetchAnswersOfExercise", exercise).Return(answers, nil)

	result, err := suite.service.FetchAnswersOfExercise(exercise)

	suite.reader.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), answers, result)
}

func (suite *ServiceSuite) TestStoreLesson() {
	lesson := &models.Lesson{Id: 1, Name: "Lesson 1"}

	suite.writer.On("StoreLesson", lesson).Return(nil)

	err := suite.service.StoreLesson(lesson)

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

func (suite *ServiceSuite) TestStoreExercise() {
	exercise := &models.Exercise{Id: 1, Question: "Exercise 1"}

	suite.writer.On("StoreExercise", exercise).Return(nil)

	err := suite.service.StoreExercise(exercise)

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

func (suite *ServiceSuite) TestStoreAnswer() {
	answer := &models.Answer{Id: 1, Type: models.Good}

	suite.writer.On("StoreAnswer", answer).Return(nil)

	err := suite.service.StoreAnswer(answer)

	suite.writer.AssertExpectations(suite.T())
	assert.NoError(suite.T(), err)
}
