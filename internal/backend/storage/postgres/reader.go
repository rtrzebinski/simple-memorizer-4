package postgres

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
)

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	return &Reader{db: db}
}

func (r *Reader) FetchLessons() (backend.Lessons, error) {
	var lessons backend.Lessons

	const query = `
		SELECT l.id, l.name, description, count(e.id) AS exercise_count
		FROM lesson l
		LEFT JOIN exercise e ON e.lesson_id = l.id
		GROUP BY l.id, l.name, description
		ORDER BY l.id, l.name, description DESC
		`

	rows, err := r.db.Query(query)
	if err != nil {
		return lessons, err
	}

	for rows.Next() {
		var lesson backend.Lesson

		err = rows.Scan(&lesson.Id, &lesson.Name, &lesson.Description, &lesson.ExerciseCount)
		if err != nil {
			return lessons, err
		}

		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (r *Reader) HydrateLesson(lesson *backend.Lesson) error {
	query := `
		SELECT name, description, count(e.id) AS exercise_count
		FROM lesson l
		LEFT JOIN exercise e ON e.lesson_id = l.id 
		WHERE l.id = $1
		GROUP BY name, description
	`

	if err := r.db.QueryRow(query, lesson.Id).Scan(&lesson.Name, &lesson.Description, &lesson.ExerciseCount); err != nil {
		return fmt.Errorf("failed to execute 'SELECT FROM lesson' query: %w", err)
	}

	return nil
}

func (r *Reader) FetchExercises(lesson backend.Lesson) (backend.Exercises, error) {
	var exercises backend.Exercises

	const query = `
		SELECT e.id, e.question, e.answer, r.id, r.type, r.created_at 
		FROM exercise e
		LEFT JOIN result r ON r.exercise_id = e.id 
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
		var resultId sql.NullInt64
		var resultType sql.NullString
		var resultCreatedAt sql.NullTime

		err = rows.Scan(&exerciseId, &exerciseQuestion, &exerciseAnswer, &resultId, &resultType, &resultCreatedAt)
		if err != nil {
			return exercises, err
		}

		result := backend.Result{}

		if resultId.Valid == true {
			numInt, err := strconv.Atoi(strconv.FormatInt(resultId.Int64, 10))
			if err != nil {
				return exercises, err
			}
			result.Id = numInt
		}
		if resultType.Valid == true {
			var ans = backend.ResultType(resultType.String)
			result.Type = ans
		}
		if resultCreatedAt.Valid == true {
			result.CreatedAt = resultCreatedAt.Time
		}

		lastIndex := len(exercises) - 1

		if lastIndex >= 0 && exercises[lastIndex].Id == exerciseId {
			// existing exercise (if exists it will always have result at this point)
			exercises[lastIndex].Results = append(exercises[lastIndex].Results, result)
		} else {
			// new exercise
			exercise := backend.Exercise{
				Id:       exerciseId,
				Question: exerciseQuestion,
				Answer:   exerciseAnswer,
			}

			// add result if exists (mind LEFT JOIN, it might be empty)
			if result.Id > 0 {
				exercise.Results = backend.Results{result}
			}

			exercises = append(exercises, exercise)
		}
	}

	return exercises, nil
}
