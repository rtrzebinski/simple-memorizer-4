package rest

import (
	"encoding/json"
	"fmt"
	myhttp "github.com/rtrzebinski/simple-memorizer-4/internal/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

type Writer struct {
	c myhttp.Client
}

func NewWriter(c myhttp.Client) *Writer {
	return &Writer{c: c}
}

func (w *Writer) StoreLesson(lesson *models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode lesson: %w", err)
	}

	_, err = w.c.Call("POST", StoreLesson, nil, body)
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

	_, err = w.c.Call("POST", DeleteLesson, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) StoreExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode exercise: %w", err)
	}

	_, err = w.c.Call("POST", StoreExercise, nil, body)
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

	_, err = w.c.Call("POST", DeleteExercise, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) IncrementBadAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode exercise: %w", err)
	}

	_, err = w.c.Call("POST", IncrementBadAnswers, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) IncrementGoodAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode exercise: %w", err)
	}

	_, err = w.c.Call("POST", IncrementGoodAnswers, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}
