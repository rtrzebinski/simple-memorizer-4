package postgres

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/stretchr/testify/assert"
)

func (suite *PostgresSuite) TestUpsertLesson_createNew() {
	db := suite.db

	w := NewWriter(db)

	lesson := models.Lesson{
		Name:        "name",
		Description: "description",
	}

	err := w.UpsertLesson(&lesson)
	assert.NoError(suite.T(), err)

	stored := fetchLatestLesson(db)

	assert.Equal(suite.T(), lesson.Name, stored.Name)
	assert.Equal(suite.T(), lesson.Description, stored.Description)
	assert.Equal(suite.T(), lesson.Id, stored.Id)
}

func (suite *PostgresSuite) TestUpsertLesson_updateExisting() {
	db := suite.db

	w := NewWriter(db)

	lesson := &Lesson{}
	createLesson(db, lesson)

	err := w.UpsertLesson(&models.Lesson{
		Id:          1,
		Name:        "newName",
		Description: "newDescription",
	})
	assert.NoError(suite.T(), err)

	stored := fetchLatestLesson(db)

	assert.Equal(suite.T(), "newName", stored.Name)
	assert.Equal(suite.T(), "newDescription", stored.Description)
}

func (suite *PostgresSuite) TestDeleteLesson() {
	db := suite.db

	w := NewWriter(db)

	createLesson(db, &Lesson{})
	stored := fetchLatestLesson(db)

	createLesson(db, &Lesson{
		Name: "another",
	})
	another := fetchLatestLesson(db)

	err := w.DeleteLesson(models.Lesson{Id: stored.Id})
	assert.NoError(suite.T(), err)

	assert.Nil(suite.T(), findLessonById(db, stored.Id))
	assert.Equal(suite.T(), "another", findLessonById(db, another.Id).Name)
}

func (suite *PostgresSuite) TestUpsertExercise_createNew() {
	db := suite.db

	w := NewWriter(db)

	lesson := &Lesson{}
	createLesson(db, lesson)

	exercise := models.Exercise{
		Lesson: &models.Lesson{
			Id: lesson.Id,
		},
		Question: "question",
		Answer:   "answer",
	}

	err := w.UpsertExercise(&exercise)
	assert.NoError(suite.T(), err)

	stored := fetchLatestExercise(db)

	assert.Equal(suite.T(), exercise.Lesson.Id, stored.LessonId)
	assert.Equal(suite.T(), exercise.Question, stored.Question)
	assert.Equal(suite.T(), exercise.Answer, stored.Answer)
	assert.Equal(suite.T(), exercise.Id, stored.Id)
}

func (suite *PostgresSuite) TestUpsertExercise_updateExisting() {
	db := suite.db

	w := NewWriter(db)

	lesson := Lesson{}
	createLesson(db, &lesson)

	exercise := Exercise{LessonId: lesson.Id}
	createExercise(db, &exercise)

	err := w.UpsertExercise(&models.Exercise{
		Id:       1,
		Question: "newQuestion",
		Answer:   "newAnswer",
	})
	assert.NoError(suite.T(), err)

	stored := fetchLatestExercise(db)

	assert.Equal(suite.T(), lesson.Id, stored.LessonId)
	assert.Equal(suite.T(), "newQuestion", stored.Question)
	assert.Equal(suite.T(), "newAnswer", stored.Answer)
}

func (suite *PostgresSuite) TestStoreExercises() {
	db := suite.db

	w := NewWriter(db)

	lesson := &Lesson{}
	createLesson(db, lesson)

	// exercise1 existing
	exercise1 := models.Exercise{
		Lesson: &models.Lesson{
			Id: lesson.Id,
		},
		Question: "question1",
		Answer:   "answer1",
	}

	// store exercise 1 to db
	createExercise(db, &Exercise{
		LessonId: exercise1.Lesson.Id,
		Question: exercise1.Question,
		Answer:   exercise1.Answer,
	})

	// exercise2 not existing
	exercise2 := models.Exercise{
		Lesson: &models.Lesson{
			Id: lesson.Id,
		},
		Question: "question2",
		Answer:   "answer2",
	}

	exercises := models.Exercises{exercise1, exercise2}

	err := w.StoreExercises(exercises)
	assert.NoError(suite.T(), err)

	ex1 := findExerciseById(db, 1)

	assert.Equal(suite.T(), exercise1.Lesson.Id, ex1.LessonId)
	assert.Equal(suite.T(), exercise1.Question, ex1.Question)
	assert.Equal(suite.T(), exercise1.Answer, ex1.Answer)

	// ID of inserted exercise will be 3, not 2,
	// this is because 'ON CONFLICT (lesson_id, question) DO NOTHING',
	// is still increasing PK auto increment value, even if nothing is inserted
	ex2 := findExerciseById(db, 3)
	assert.Equal(suite.T(), exercise2.Lesson.Id, ex2.LessonId)
	assert.Equal(suite.T(), exercise2.Question, ex2.Question)
	assert.Equal(suite.T(), exercise2.Answer, ex2.Answer)
}

func (suite *PostgresSuite) TestDeleteExercise() {
	db := suite.db

	w := NewWriter(db)

	lesson := &Lesson{
		ExerciseCount: 1,
	}
	createLesson(db, lesson)

	createExercise(db, &Exercise{
		LessonId: lesson.Id,
	})
	stored := fetchLatestExercise(db)

	createExercise(db, &Exercise{
		Question: "another",
	})
	another := fetchLatestExercise(db)

	err := w.DeleteExercise(models.Exercise{Id: stored.Id})
	assert.NoError(suite.T(), err)

	assert.Nil(suite.T(), findExerciseById(db, stored.Id))
	assert.Equal(suite.T(), "another", findExerciseById(db, another.Id).Question)

	lesson = findLessonById(db, stored.LessonId)
	assert.Equal(suite.T(), 0, lesson.ExerciseCount)
}

func (suite *PostgresSuite) TestStoreAnswer() {
	db := suite.db

	w := NewWriter(db)

	exercise := &Exercise{}
	createExercise(db, exercise)

	answer := models.Result{
		Type: models.Good,
		Exercise: &models.Exercise{
			Id: exercise.Id,
		},
	}

	err := w.StoreResult(&answer)
	assert.NoError(suite.T(), err)

	stored := fetchLatestResult(db)

	assert.Equal(suite.T(), answer.Type, stored.Type)
	assert.Equal(suite.T(), answer.Exercise.Id, stored.ExerciseId)
}
