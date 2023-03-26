package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"net/http"
)

// A Home component
type Home struct {
	app.Compo
	api *frontend.ApiClient

	isAnswerVisible bool
	isNextPreloaded bool

	// currently visible exercise
	question    string
	answer      string
	goodAnswers int
	badAnswers  int
	exerciseId  int

	// next exercise preloaded
	nextQuestion    string
	nextAnswer      string
	nextGoodAnswers int
	nextBadAnswers  int
	nextExerciseId  int
}

// The OnMount method is run once component is mounted
func (h *Home) OnMount(ctx app.Context) {
	url := app.Window().URL()
	h.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
	h.handleNextExercise()

	// Bind actions to keyboard shortcuts
	app.Window().AddEventListener("keyup", func(ctx app.Context, e app.Event) {
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
		case "KeyN":
			h.handleNextExercise()
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

func (h *Home) handleNextExercise() {
	h.isAnswerVisible = false

	if h.isNextPreloaded == false {
		exercise := h.reload()
		h.exerciseId = exercise.Id
		h.question = exercise.Question
		h.answer = exercise.Answer
		h.goodAnswers = exercise.GoodAnswers
		h.badAnswers = exercise.BadAnswers
		app.Log("displayed initial exercise")
	} else {
		h.exerciseId = h.nextExerciseId
		h.question = h.nextQuestion
		h.answer = h.nextAnswer
		h.goodAnswers = h.nextGoodAnswers
		h.badAnswers = h.nextBadAnswers
		app.Log("displayed preloaded exercise")
		h.isNextPreloaded = false
	}

	go func() {
		exercise := h.reload()
		h.nextExerciseId = exercise.Id
		h.nextQuestion = exercise.Question
		h.nextAnswer = exercise.Answer
		h.nextGoodAnswers = exercise.GoodAnswers
		h.nextBadAnswers = exercise.BadAnswers
		app.Log("preloaded")
		h.isNextPreloaded = true
	}()
}

func (h *Home) reload() models.Exercise {
	exercise, err := h.api.FetchNextExercise()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch next exercise: %w", err))
	}

	if exercise.Id == h.exerciseId {
		return h.reload()
	}

	return exercise
}

func (h *Home) handleViewAnswer() {
	h.isAnswerVisible = true
}

func (h *Home) handleGoodAnswer() {
	// make a copy of h.exerciseId to prevent h.exerciseId being updated before increment was completed
	toIncrement := h.exerciseId
	go func() {
		// increment in the background
		if err := h.api.IncrementGoodAnswers(toIncrement); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	h.handleNextExercise()
}

func (h *Home) handleBadAnswer() {
	// make a copy of h.exerciseId to prevent h.exerciseId being updated before increment was completed
	toIncrement := h.exerciseId
	go func() {
		// increment in the background
		if err := h.api.IncrementBadAnswers(toIncrement); err != nil {
			app.Log(fmt.Errorf("failed to increment bad answers: %w", err))
		}
	}()
	h.handleNextExercise()
}
