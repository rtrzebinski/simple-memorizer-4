package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"net/http"
)

// A Home component
type Home struct {
	app.Compo
	api *frontend.ApiClient

	question        string
	answer          string
	isAnswerVisible bool
	goodAnswers     int
	badAnswers      int
	exerciseId      int
}

// The OnMount method is run once component is mounted
func (h *Home) OnMount(ctx app.Context) {
	url := app.Window().URL()
	h.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
	h.handleNextExercise()

	// Bind actions to keyboard shortcuts
	app.Window().AddEventListener("keydown", func(ctx app.Context, e app.Event) {
		switch e.Get("code").String() {
		case "Space":
			if h.isAnswerVisible == true {
				h.handleNextExercise()
			} else {
				h.handleViewAnswer()
			}
		case "KeyV":
			h.handleViewAnswer()
		case "KeyG":
			h.handleGoodAnswer()
		case "KeyB":
			h.handleBadAnswer()
		}
	})
}

// The Render method is where the component appearance is defined.
func (h *Home) Render() app.UI {
	return app.Div().Body(
		app.P().Body(
			app.Button().
				Text("Next exercise").
				OnClick(func(ctx app.Context, e app.Event) {
					h.handleNextExercise()
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("View answer").
				OnClick(func(ctx app.Context, e app.Event) {
					h.handleViewAnswer()
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("Good answer").
				OnClick(func(ctx app.Context, e app.Event) {
					h.handleGoodAnswer()
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("Bad answer").
				OnClick(func(ctx app.Context, e app.Event) {
					h.handleBadAnswer()
				}).
				Style("margin-right", "10px"),
		),
		app.H2().Body(
			app.Text("What is the capital of "),
			app.If(h.question != "",
				app.Text(h.question),
			).Else(
				app.Text(""),
			),
			app.Text("?"),
		),
		app.H2().Body(
			app.If(h.answer != "",
				app.Text(h.answer),
			).Else(
				app.Text(""),
			),
		).Hidden(!h.isAnswerVisible),
		app.P().Body(
			app.Text("Good answers: "),
			app.Text(h.goodAnswers),
		),
		app.P().Body(
			app.Text("Bad answers: "),
			app.Text(h.badAnswers),
		),
	)
}

// handleNextExercise fetch random exercises, isAnswerVisible question, hide answer
func (h *Home) handleNextExercise() {
	exercise, err := h.api.FetchRandomExercise()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch random exercise: %w", err))
	}
	h.exerciseId = exercise.Id
	h.isAnswerVisible = false
	h.question = exercise.Question
	h.answer = exercise.Answer
	h.goodAnswers = exercise.GoodAnswers
	h.badAnswers = exercise.BadAnswers
}

func (h *Home) handleViewAnswer() {
	h.isAnswerVisible = true
}

func (h *Home) handleGoodAnswer() {
	go func() {
		err := h.api.IncrementGoodAnswers(h.exerciseId)
		if err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	h.goodAnswers++
	h.handleNextExercise()
}

func (h *Home) handleBadAnswer() {
	go func() {
		err := h.api.IncrementBadAnswers(h.exerciseId)
		if err != nil {
			app.Log(fmt.Errorf("failed to increment bad answers: %w", err))
		}
	}()
	h.badAnswers++
	h.handleNextExercise()
}
