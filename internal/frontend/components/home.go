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
	exerciseId  int
	question    string
	answer      string
	goodAnswers int
	badAnswers  int

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
	h.bindKeys()
	h.bindSwipes()
}

// The Render method is where the component appearance is defined.
func (h *Home) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
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
		app.P().Body(
			app.Button().
				Text("⇧ View answer").
				OnClick(func(ctx app.Context, e app.Event) {
					h.handleViewAnswer()
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
			app.Button().
				Text("Next exercise ⇩").
				OnClick(func(ctx app.Context, e app.Event) {
					// only allow if next exercise was preloaded (to avoid double clicks)
					if h.isNextPreloaded == true {
						h.handleNextExercise()
					}
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
		),
		app.P().Body(
			app.Button().
				Text("⇦ Bad answer").
				OnClick(func(ctx app.Context, e app.Event) {
					// only allow if next exercise was preloaded (to avoid double clicks)
					if h.isNextPreloaded == true {
						h.handleBadAnswer()
					}
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
			app.Button().
				Text("Good answer ⇨").
				OnClick(func(ctx app.Context, e app.Event) {
					// only allow if next exercise was preloaded (to avoid double clicks)
					if h.isNextPreloaded == true {
						h.handleGoodAnswer()
					}
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
		),
	)
}

func (h *Home) bindKeys() {
	app.Window().AddEventListener("keyup", func(ctx app.Context, e app.Event) {
		// bind actions to keyboard shortcuts
		switch e.Get("code").String() {
		case "Space":
			if h.isAnswerVisible == true {
				// only allow if next exercise was preloaded (to avoid double clicks)
				if h.isNextPreloaded == true {
					h.handleNextExercise()
				}
			} else {
				h.handleViewAnswer()
			}
		case "KeyV", "ArrowUp":
			h.handleViewAnswer()
		case "KeyG", "ArrowRight":
			// only allow if next exercise was preloaded (to avoid double clicks)
			if h.isNextPreloaded == true {
				h.handleGoodAnswer()
			}
		case "KeyB", "ArrowLeft":
			// only allow if next exercise was preloaded (to avoid double clicks)
			if h.isNextPreloaded == true {
				h.handleBadAnswer()
			}
		case "KeyN", "ArrowDown":
			// only allow if next exercise was preloaded (to avoid double clicks)
			if h.isNextPreloaded == true {
				h.handleNextExercise()
			}
		}
	})
}

func (h *Home) bindSwipes() {
	app.Window().AddEventListener("swiped-left", func(ctx app.Context, e app.Event) {
		// only allow if next exercise was preloaded (to avoid double clicks)
		if h.isNextPreloaded == true {
			h.handleBadAnswer()
		}
	})
	app.Window().AddEventListener("swiped-right", func(ctx app.Context, e app.Event) {
		// only allow if next exercise was preloaded (to avoid double clicks)
		if h.isNextPreloaded == true {
			h.handleGoodAnswer()
		}
	})
	app.Window().AddEventListener("swiped-up", func(ctx app.Context, e app.Event) {
		// only allow if next exercise was preloaded (to avoid double clicks)
		if h.isNextPreloaded == true {
			h.handleNextExercise()
		}
	})
	app.Window().AddEventListener("swiped-down", func(ctx app.Context, e app.Event) {
		h.handleViewAnswer()
	})
}

func (h *Home) handleNextExercise() {
	h.isAnswerVisible = false

	if h.isNextPreloaded == false {
		exercise := h.fetchNext()
		h.exerciseId = exercise.Id
		h.question = exercise.Question
		h.answer = exercise.Answer
		h.goodAnswers = exercise.GoodAnswers
		h.badAnswers = exercise.BadAnswers
		app.Log("displayed initial exercise")
	} else {
		h.isNextPreloaded = false
		h.exerciseId = h.nextExerciseId
		h.question = h.nextQuestion
		h.answer = h.nextAnswer
		h.goodAnswers = h.nextGoodAnswers
		h.badAnswers = h.nextBadAnswers
		app.Log("displayed preloaded exercise")
	}

	go func() {
		exercise := h.fetchNext()
		h.nextExerciseId = exercise.Id
		h.nextQuestion = exercise.Question
		h.nextAnswer = exercise.Answer
		h.nextGoodAnswers = exercise.GoodAnswers
		h.nextBadAnswers = exercise.BadAnswers
		h.isNextPreloaded = true
		app.Log("preloaded")
	}()
}

func (h *Home) fetchNext() models.Exercise {
	exercise, err := h.api.FetchNextExercise()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch next exercise: %w", err))
	}

	// dummy way of avoiding duplicates todo move to the API
	if exercise.Id == h.exerciseId {
		return h.fetchNext()
	}

	return exercise
}

func (h *Home) handleViewAnswer() {
	h.isAnswerVisible = true
}

func (h *Home) handleGoodAnswer() {
	exercise := models.Exercise{Id: h.exerciseId}
	go func() {
		// increment in the background
		if err := h.api.IncrementGoodAnswers(exercise); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	h.handleNextExercise()
}

func (h *Home) handleBadAnswer() {
	exercise := models.Exercise{Id: h.exerciseId}
	go func() {
		// increment in the background
		if err := h.api.IncrementBadAnswers(exercise); err != nil {
			app.Log(fmt.Errorf("failed to increment bad answers: %w", err))
		}
	}()
	h.handleNextExercise()
}
