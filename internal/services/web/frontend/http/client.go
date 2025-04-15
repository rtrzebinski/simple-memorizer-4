package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
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
func (s *Client) FetchLessons(ctx context.Context, accessToken string) ([]frontend.Lesson, error) {
	var lessons []frontend.Lesson

	respBody, err := s.c.Call(ctx, "GET", http.FetchLessons, nil, nil, accessToken)
	if err != nil {
		return lessons, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &lessons); err != nil {
		return lessons, fmt.Errorf("failed to decode lessons: %w", err)
	}

	return lessons, nil
}

// HydrateLesson hydrates a lesson
func (s *Client) HydrateLesson(ctx context.Context, lesson *frontend.Lesson, accessToken string) error {
	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := s.c.Call(ctx, "GET", http.HydrateLesson, params, nil, accessToken)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, lesson); err != nil {
		return fmt.Errorf("failed to decode lesson: %w", err)
	}

	return nil
}

// FetchExercises fetches exercises of a lesson
func (s *Client) FetchExercises(ctx context.Context, lesson frontend.Lesson, accessToken string) ([]frontend.Exercise, error) {
	var exercises []frontend.Exercise

	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := s.c.Call(ctx, "GET", http.FetchExercises, params, nil, accessToken)
	if err != nil {
		return exercises, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &exercises); err != nil {
		return exercises, fmt.Errorf("failed to decode exercises: %w", err)
	}

	return exercises, nil
}

// UpsertLesson upserts a lesson
func (s *Client) UpsertLesson(ctx context.Context, lesson frontend.Lesson, accessToken string) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode lesson: %w", err)
	}

	_, err = s.c.Call(ctx, "POST", http.UpsertLesson, nil, body, accessToken)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// DeleteLesson deletes a lesson
func (s *Client) DeleteLesson(ctx context.Context, lesson frontend.Lesson, accessToken string) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode lesson: %w", err)
	}

	_, err = s.c.Call(ctx, "POST", http.DeleteLesson, nil, body, accessToken)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// UpsertExercise upserts an exercise
func (s *Client) UpsertExercise(ctx context.Context, exercise frontend.Exercise, accessToken string) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode exercise: %w", err)
	}

	_, err = s.c.Call(ctx, "POST", http.UpsertExercise, nil, body, accessToken)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// StoreExercises stores exercises
func (s *Client) StoreExercises(ctx context.Context, exercises []frontend.Exercise, accessToken string) error {
	body, err := json.Marshal(exercises)
	if err != nil {
		return fmt.Errorf("failed to encode exercises: %w", err)
	}

	_, err = s.c.Call(ctx, "POST", http.StoreExercises, nil, body, accessToken)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// DeleteExercise deletes an exercise
func (s *Client) DeleteExercise(ctx context.Context, exercise frontend.Exercise, accessToken string) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode exercise: %w", err)
	}

	_, err = s.c.Call(ctx, "POST", http.DeleteExercise, nil, body, accessToken)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

// StoreResult stores a result
func (s *Client) StoreResult(ctx context.Context, result frontend.Result, accessToken string) error {
	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	_, err = s.c.Call(ctx, "POST", http.StoreResult, nil, body, accessToken)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (s *Client) AuthRegister(ctx context.Context, req frontend.RegisterRequest) (frontend.RegisterResponse, error) {
	res := frontend.RegisterResponse{}

	body, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("failed to encode register request: %w", err)
	}

	respBody, err := s.c.Call(ctx, "POST", http.AuthRegister, nil, body, "")
	if err != nil {
		return res, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &res); err != nil {
		return res, fmt.Errorf("failed to decode register response: %w", err)
	}

	return res, nil
}

func (s *Client) AuthSignIn(ctx context.Context, req frontend.SignInRequest) (frontend.SignInResponse, error) {
	res := frontend.SignInResponse{}

	body, err := json.Marshal(req)
	if err != nil {
		return res, fmt.Errorf("failed to encode sign in request: %w", err)
	}

	respBody, err := s.c.Call(ctx, "POST", http.AuthSignIn, nil, body, "")
	if err != nil {
		return res, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &res); err != nil {
		return res, fmt.Errorf("failed to decode sign in response: %w", err)
	}

	return res, nil
}
