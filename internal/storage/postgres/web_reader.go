package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type WebReader struct {
	db *sql.DB
}

func NewWebReader(db *sql.DB) *WebReader {
	return &WebReader{db: db}
}

func (r *WebReader) FetchLessons(ctx context.Context, userID string) (backend.Lessons, error) {
	slog.Debug("WebReader FetchLessons", "userID", userID)

	var lessons backend.Lessons

	const query = `
		SELECT l.id, l.name, description, count(e.id) AS exercise_count
		FROM lesson l
		LEFT JOIN exercise e ON e.lesson_id = l.id
		WHERE l.user_id = $1
		GROUP BY l.id, l.name, description
		ORDER BY l.id, l.name, description DESC
		`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return lessons, fmt.Errorf("failed to fetch lessons: %w", err)
	}
	defer rows.Close()

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

func (r *WebReader) HydrateLesson(ctx context.Context, userID string, lesson *backend.Lesson) error {
	slog.Debug("WebReader HydrateLesson", "userID", userID)

	query := `
		SELECT name, description, count(e.id) AS exercise_count
		FROM lesson l
		LEFT JOIN exercise e ON e.lesson_id = l.id 
		WHERE l.id = $1 AND l.user_id = $2
		GROUP BY name, description
	`

	err := r.db.QueryRowContext(ctx, query, lesson.Id, userID).Scan(&lesson.Name, &lesson.Description, &lesson.ExerciseCount)
	if err != nil {
		return fmt.Errorf("failed to hydrate lesson: %w", err)
	}

	return nil
}

func (r *WebReader) FetchExercises(ctx context.Context, userID string, lesson backend.Lesson, oldestExerciseID int) (backend.Exercises, error) {
	slog.Debug("WebReader FetchExercises", "userID", userID)

	const query = `
		SELECT e.id, e.question, e.answer,
       	e.bad_answers, e.bad_answers_today, e.latest_bad_answer, e.latest_bad_answer_was_today,
       	e.good_answers, e.good_answers_today, e.latest_good_answer, e.latest_good_answer_was_today
		FROM exercise e
		JOIN lesson l ON e.lesson_id = l.id
		WHERE e.lesson_id = $1
  		AND e.id >= $2
  		AND l.user_id = $3
		ORDER BY e.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, lesson.Id, oldestExerciseID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exercises: %w", err)
	}
	defer rows.Close()

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
