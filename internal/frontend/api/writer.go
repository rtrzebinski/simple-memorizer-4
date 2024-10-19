package api

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
)

type Writer struct {
	c Caller
}

func NewWriter(c Caller) *Writer {
	return &Writer{c: c}
}

func (w *Writer) UpsertLesson(lesson *models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode lesson: %w", err)
	}

	_, err = w.c.Call("POST", routes.UpsertLesson, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) DeleteLesson(lesson models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode lesson: %w", err)
	}

	_, err = w.c.Call("POST", routes.DeleteLesson, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) UpsertExercise(exercise *models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode exercise: %w", err)
	}

	_, err = w.c.Call("POST", routes.UpsertExercise, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) StoreExercises(exercises models.Exercises) error {
	body, err := json.Marshal(exercises)
	if err != nil {
		return fmt.Errorf("failed to encode exercises: %w", err)
	}

	_, err = w.c.Call("POST", routes.StoreExercises, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) DeleteExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode exercise: %w", err)
	}

	_, err = w.c.Call("POST", routes.DeleteExercise, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) StoreResult(result *models.Result) error {
	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	_, err = w.c.Call("POST", routes.StoreResult, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}
