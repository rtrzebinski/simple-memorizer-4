package postgres

import (
	"database/sql"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"strconv"
)

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	return &Reader{db: db}
}

func (r *Reader) FetchAllLessons() (models.Lessons, error) {
	var lessons models.Lessons

	const query = `
		SELECT l.id, l.name, l.description, l.exercise_count
		FROM lesson l
		ORDER BY l.id DESC
		`

	rows, err := r.db.Query(query)
	if err != nil {
		return lessons, err
	}

	for rows.Next() {
		var lesson models.Lesson

		err = rows.Scan(&lesson.Id, &lesson.Name, &lesson.Description, &lesson.ExerciseCount)
		if err != nil {
			return lessons, err
		}

		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (r *Reader) HydrateLesson(lesson *models.Lesson) error {
	query := `SELECT name, description, exercise_count FROM lesson WHERE id = $1;`

	if err := r.db.QueryRow(query, lesson.Id).Scan(&lesson.Name, &lesson.Description, &lesson.ExerciseCount); err != nil {
		return fmt.Errorf("failed to execute 'SELECT FROM lesson' query: %w", err)
	}

	return nil
}

func (r *Reader) FetchExercisesOfLesson(lesson models.Lesson) (models.Exercises, error) {
	var exercises models.Exercises

	const query = `
		SELECT e.id, e.question, e.answer, a.id, a.type, a.created_at 
		FROM exercise e
		LEFT JOIN answer a ON a.exercise_id = e.id 
		WHERE lesson_id = $1
		ORDER BY e.id DESC
		`

	rows, err := r.db.Query(query, lesson.Id)
	if err != nil {
		return exercises, err
	}

	for rows.Next() {
		var exerciseId int
		var exerciseQuestion string
		var exerciseAnswer string
		var answerId sql.NullInt64
		var answerType sql.NullString
		var answerCreatedAt sql.NullTime

		err = rows.Scan(&exerciseId, &exerciseQuestion, &exerciseAnswer, &answerId, &answerType, &answerCreatedAt)
		if err != nil {
			return exercises, err
		}

		answer := models.Answer{}

		if answerId.Valid == true {
			numInt, err := strconv.Atoi(strconv.FormatInt(answerId.Int64, 10))
			if err != nil {
				return exercises, err
			}
			answer.Id = numInt
		}
		if answerType.Valid == true {
			var ans = models.AnswerType(answerType.String)
			answer.Type = ans
		}
		if answerCreatedAt.Valid == true {
			answer.CreatedAt = answerCreatedAt.Time
		}

		lastIndex := len(exercises) - 1

		if lastIndex >= 0 && exercises[lastIndex].Id == exerciseId {
			// existing exercise (if exists it will always have answer at this point)
			exercises[lastIndex].Answers = append(exercises[lastIndex].Answers, answer)
		} else {
			// new exercise
			exercise := models.Exercise{
				Id:       exerciseId,
				Question: exerciseQuestion,
				Answer:   exerciseAnswer,
			}

			// add answer if exists (mind LEFT JOIN, it might be empty)
			if answer.Id > 0 {
				exercise.Answers = models.Answers{answer}
			}

			exercises = append(exercises, exercise)
		}
	}

	return exercises, nil
}

func (r *Reader) FetchAnswersOfExercise(exercise models.Exercise) (models.Answers, error) {
	answers := models.Answers{}

	const query = `
		SELECT a.id, a.type, a.created_at
		FROM answer a
		WHERE a.exercise_id = $1
		`

	rows, err := r.db.Query(query, exercise.Id)
	if err != nil {
		return answers, err
	}

	for rows.Next() {
		var answer models.Answer

		err = rows.Scan(&answer.Id, &answer.Type, &answer.CreatedAt)
		if err != nil {
			return answers, err
		}

		answers = append(answers, answer)
	}

	return answers, nil
}
