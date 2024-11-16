package api

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"strconv"
)

// Client is a client for the API
type Client struct {
	c Caller
}

// NewClient creates a new Client
func NewClient(c Caller) *Client {
	return &Client{c: c}
}

// FetchLessons fetches all lessons
func (s *Client) FetchLessons() ([]frontend.Lesson, error) {
	var lessons []frontend.Lesson

	respBody, err := s.c.Call("GET", server.FetchLessons, nil, nil)
	if err != nil {
		return lessons, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &lessons); err != nil {
		return lessons, fmt.Errorf("failed to decode lessons: %w", err)
	}

	return lessons, nil
}

// HydrateLesson hydrates a lesson
func (s *Client) HydrateLesson(lesson *frontend.Lesson) error {
	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := s.c.Call("GET", server.HydrateLesson, params, nil)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, lesson); err != nil {
		return fmt.Errorf("failed to decode lesson: %w", err)
	}

	return nil
}

// FetchExercises fetches exercises of a lesson
func (s *Client) FetchExercises(lesson frontend.Lesson) ([]frontend.Exercise, error) {
	var exercises []frontend.Exercise

	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := s.c.Call("GET", server.FetchExercises, params, nil)
	if err != nil {
		return exercises, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &exercises); err != nil {
		return exercises, fmt.Errorf("failed to decode exercises: %w", err)
	}

	return exercises, nil
}

// UpsertLesson upserts a lesson
func (s *Client) UpsertLesson(lesson frontend.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode lesson: %w", err)
	}

	_, err = s.c.Call("POST", server.UpsertLesson, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// DeleteLesson deletes a lesson
func (s *Client) DeleteLesson(lesson frontend.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode lesson: %w", err)
	}

	_, err = s.c.Call("POST", server.DeleteLesson, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// UpsertExercise upserts an exercise
func (s *Client) UpsertExercise(exercise frontend.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode exercise: %w", err)
	}

	_, err = s.c.Call("POST", server.UpsertExercise, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// StoreExercises stores exercises
func (s *Client) StoreExercises(exercises []frontend.Exercise) error {
	body, err := json.Marshal(exercises)
	if err != nil {
		return fmt.Errorf("failed to encode exercises: %w", err)
	}

	_, err = s.c.Call("POST", server.StoreExercises, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// DeleteExercise deletes an exercise
func (s *Client) DeleteExercise(exercise frontend.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode exercise: %w", err)
	}

	_, err = s.c.Call("POST", server.DeleteExercise, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// StoreResult stores a result
func (s *Client) StoreResult(result frontend.Result) error {
	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	_, err = s.c.Call("POST", server.StoreResult, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}
