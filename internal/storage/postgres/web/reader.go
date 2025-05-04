package web

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	return &Reader{db: db}
}

func (r *Reader) FetchLessons(ctx context.Context, userID string) (backend.Lessons, error) {
	slog.Debug("Reader FetchLessons", "userID", userID)

	var lessons backend.Lessons

	const query = `
		SELECT l.id, l.name, description, count(e.id) AS exercise_count
		FROM lesson l
		LEFT JOIN exercise e ON e.lesson_id = l.id
		GROUP BY l.id, l.name, description
		ORDER BY l.id, l.name, description DESC
		`

	rows, err := r.db.QueryContext(ctx, query)
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

func (r *Reader) HydrateLesson(ctx context.Context, lesson *backend.Lesson, userID string) error {
	slog.Debug("Reader HydrateLesson", "userID", userID)

	query := `
		SELECT name, description, count(e.id) AS exercise_count
		FROM lesson l
		LEFT JOIN exercise e ON e.lesson_id = l.id 
		WHERE l.id = $1
		GROUP BY name, description
	`

	if err := r.db.QueryRowContext(ctx, query, lesson.Id).Scan(&lesson.Name, &lesson.Description, &lesson.ExerciseCount); err != nil {
		return fmt.Errorf("failed to execute 'SELECT FROM lesson' query: %w", err)
	}

	return nil
}

func (r *Reader) FetchExercises(ctx context.Context, lesson backend.Lesson, userID string) (backend.Exercises, error) {
	slog.Debug("Reader FetchExercises", "userID", userID)

	const query = `
SELECT e.id, e.question, e.answer,
       e.bad_answers, e.bad_answers_today, e.latest_bad_answer, e.latest_bad_answer_was_today,
       e.good_answers, e.good_answers_today, e.latest_good_answer, e.latest_good_answer_was_today
FROM exercise e
WHERE lesson_id = $1
ORDER BY e.id DESC
`

	rows, err := r.db.QueryContext(ctx, query, lesson.Id)
	if err != nil {
		return nil, err
	}

	var exercises backend.Exercises

	for rows.Next() {
		var exercise backend.Exercise

		err = rows.Scan(&exercise.Id, &exercise.Question, &exercise.Answer, &exercise.BadAnswers, &exercise.BadAnswersToday,
			&exercise.LatestBadAnswer, &exercise.LatestBadAnswerWasToday, &exercise.GoodAnswers, &exercise.GoodAnswersToday,
			&exercise.LatestGoodAnswer, &exercise.LatestGoodAnswerWasToday)
		if err != nil {
			return nil, err
		}

		exercises = append(exercises, exercise)
	}

	return exercises, nil
}
