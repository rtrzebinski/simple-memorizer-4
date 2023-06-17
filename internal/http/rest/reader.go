package rest

import (
	"encoding/json"
	"fmt"
	myhttp "github.com/rtrzebinski/simple-memorizer-4/internal/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"strconv"
)

type Reader struct {
	c myhttp.Client
}

func NewReader(c myhttp.Client) *Reader {
	return &Reader{c: c}
}

func (r *Reader) FetchAllLessons() (models.Lessons, error) {
	var lessons models.Lessons

	respBody, err := r.c.Call("GET", FetchAllLessons, nil, nil)
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

	respBody, err := r.c.Call("GET", HydrateLesson, params, nil)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, lesson); err != nil {
		return fmt.Errorf("failed to decode lesson: %w", err)
	}

	return nil
}

func (r *Reader) FetchExercisesOfLesson(lesson models.Lesson) (models.Exercises, error) {
	var exercises models.Exercises

	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := r.c.Call("GET", FetchExercisesOfLesson, params, nil)
	if err != nil {
		return exercises, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &exercises); err != nil {
		return exercises, fmt.Errorf("failed to decode exercises: %w", err)
	}

	return exercises, nil
}

func (r *Reader) FetchAnswersOfExercise(exercise models.Exercise) (models.Answers, error) {
	var answers models.Answers

	var params = map[string]string{
		"exercise_id": strconv.Itoa(exercise.Id),
	}

	respBody, err := r.c.Call("GET", FetchAnswersOfExercise, params, nil)
	if err != nil {
		return answers, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &answers); err != nil {
		return answers, fmt.Errorf("failed to decode answers: %w", err)
	}

	return answers, nil
}
