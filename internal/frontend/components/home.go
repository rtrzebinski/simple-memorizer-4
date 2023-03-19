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

	question    string
	answer      string
	showAnswer  bool
	goodAnswers int
	badAnswers  int
	exerciseId  int
}

// The OnMount method is run once component is mounted
func (h *Home) OnMount(ctx app.Context) {
	url := app.Window().URL()
	h.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
	h.nextExercise()
}

// The Render method is where the component appearance is defined.
func (h *Home) Render() app.UI {
	return app.Div().Body(
		app.P().Body(
			app.Button().
				Text("Next exercise").
				OnClick(func(ctx app.Context, e app.Event) {
					h.nextExercise()
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("Show Answer").
				OnClick(func(ctx app.Context, e app.Event) {
					h.showAnswer = true
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("Good Answer").
				OnClick(func(ctx app.Context, e app.Event) {
					go func() {
						err := h.api.IncrementGoodAnswers(h.exerciseId)
						if err != nil {
							app.Log(fmt.Errorf("failed to increment good answers: %w", err))
						}
					}()
					h.goodAnswers++
					h.nextExercise()
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("Bad Answer").
				OnClick(func(ctx app.Context, e app.Event) {
					go func() {
						err := h.api.IncrementBadAnswers(h.exerciseId)
						if err != nil {
							app.Log(fmt.Errorf("failed to increment bad answers: %w", err))
						}
					}()
					h.badAnswers++
					h.nextExercise()
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
		).Hidden(!h.showAnswer),
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

// nextExercise fetch random exercises, showAnswer question, hide answer
func (h *Home) nextExercise() {
	exercise, err := h.api.FetchRandomExercise()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch random exercise: %w", err))
	}
	h.exerciseId = exercise.Id
	h.showAnswer = false
	h.question = exercise.Question
	h.answer = exercise.Answer
	h.goodAnswers = exercise.GoodAnswers
	h.badAnswers = exercise.BadAnswers
}
