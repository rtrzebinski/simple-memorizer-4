package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-go/internal/api"
)

// A Home component
type Home struct {
	app.Compo
	client api.Client

	question    string
	answer      string
	showAnswer  bool
	goodAnswers int
	badAnswers  int
}

// NewHome component constructor
func NewHome(client api.Client) *Home {
	return &Home{
		client: client,
	}
}

// fetchExercise fetch random exercises, showAnswer question, hide answer
func (h *Home) fetchExercise() {
	h.showAnswer = false
	exercise := h.client.GetRandomExercise()
	h.question = exercise.Question
	h.answer = exercise.Answer
}

// The OnMount method is run once component is mounted
func (h *Home) OnMount(ctx app.Context) {
	// host can only be read from the Window once component is mounted
	h.client.SetHost(app.Window().URL().Host)
	// scheme can only be read from the Window once component is mounted
	h.client.SetScheme(app.Window().URL().Scheme)
	// knowing host, fetch and display the initial exercise
	h.fetchExercise()
}

// The Render method is where the component appearance is defined.
func (h *Home) Render() app.UI {
	return app.Div().Body(
		app.P().Body(
			app.Button().
				Text("Next exercise").
				OnClick(func(ctx app.Context, e app.Event) {
					h.fetchExercise()
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
					h.goodAnswers++
					h.fetchExercise()
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("Bad Answer").
				OnClick(func(ctx app.Context, e app.Event) {
					h.badAnswers++
					h.fetchExercise()
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
