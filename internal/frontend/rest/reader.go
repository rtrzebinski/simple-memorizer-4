package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"strconv"
)

type Reader struct {
	c frontend.Client
}

func NewReader(c frontend.Client) *Reader {
	return &Reader{c: c}
}

func (r *Reader) FetchLessons() (models.Lessons, error) {
	var lessons models.Lessons

	respBody, err := r.c.Call("GET", routes.FetchLessons, nil, nil)
	if err != nil {
		return lessons, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &lessons); err != nil {
		return lessons, fmt.Errorf("failed to decode lessons: %w", err)
	}

	return lessons, nil
}

func (r *Reader) HydrateLesson(lesson *models.Lesson) error {
	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := r.c.Call("GET", routes.HydrateLesson, params, nil)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, lesson); err != nil {
		return fmt.Errorf("failed to decode lesson: %w", err)
	}

	return nil
}

func (r *Reader) FetchExercises(lesson models.Lesson) (models.Exercises, error) {
	var exercises models.Exercises

	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := r.c.Call("GET", routes.FetchExercises, params, nil)
	if err != nil {
		return exercises, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &exercises); err != nil {
		return exercises, fmt.Errorf("failed to decode exercises: %w", err)
	}

	return exercises, nil
}
