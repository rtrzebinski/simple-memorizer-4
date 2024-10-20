package postgres

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/stretchr/testify/assert"
)

func (suite *PostgresSuite) TestFetchLessons() {
	db := suite.db

	r := NewReader(db)

	lesson := &Lesson{}
	createLesson(db, lesson)

	createExercise(db, &Exercise{LessonId: lesson.Id})

	res, err := r.FetchLessons()

	assert.NoError(suite.T(), err)
	assert.IsType(suite.T(), models.Lessons{}, res)
	assert.Len(suite.T(), res, 1)
	assert.Equal(suite.T(), lesson.Id, res[0].Id)
	assert.Equal(suite.T(), lesson.Name, res[0].Name)
	assert.Equal(suite.T(), lesson.Description, res[0].Description)
	assert.Equal(suite.T(), 1, res[0].ExerciseCount)
}

func (suite *PostgresSuite) TestHydrateLesson() {
	db := suite.db

	r := NewReader(db)

	l := &Lesson{}
	createLesson(db, l)

	lesson := &models.Lesson{
		Id: l.Id,
	}

	err := r.HydrateLesson(lesson)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), l.Name, lesson.Name)
	assert.Equal(suite.T(), l.Description, lesson.Description)
	assert.Equal(suite.T(), 0, lesson.ExerciseCount)

	createExercise(db, &Exercise{LessonId: l.Id})
	createExercise(db, &Exercise{LessonId: l.Id})

	err = r.HydrateLesson(lesson)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), l.Name, lesson.Name)
	assert.Equal(suite.T(), l.Description, lesson.Description)
	assert.Equal(suite.T(), 2, lesson.ExerciseCount)
}

func (suite *PostgresSuite) TestFetchExercises() {
	db := suite.db

	r := NewReader(db)

	exercise := &Exercise{}
	createExercise(db, exercise)

	// to check of exercise without results will also be fetched
	exercise2 := &Exercise{LessonId: exercise.LessonId}
	createExercise(db, exercise2)

	// belongs to exercise, to be included
	result1 := &Result{ExerciseId: exercise.Id}
	createResult(db, result1)

	// does not belong to exercise, to be excluded
	result2 := &Result{}
	createResult(db, result2)

	res, err := r.FetchExercises(models.Lesson{Id: exercise.LessonId})

	assert.NoError(suite.T(), err)
	assert.IsType(suite.T(), models.Exercises{}, res)
	assert.Len(suite.T(), res, 2)
	assert.Equal(suite.T(), exercise.Id, res[1].Id)
	assert.Equal(suite.T(), exercise.Question, res[1].Question)
	assert.Equal(suite.T(), exercise.Answer, res[1].Answer)
	assert.Len(suite.T(), res[1].Results, 1)
	assert.Empty(suite.T(), res[0].Results)
	assert.Equal(suite.T(), result1.Id, res[1].Results[0].Id)
	assert.Equal(suite.T(), result1.Type, res[1].Results[0].Type)
	assert.Equal(suite.T(), result1.CreatedAt, res[1].Results[0].CreatedAt)
}
