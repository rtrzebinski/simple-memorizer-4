package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"net/http"
)

var pathHome = "/"

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
func (c *Home) OnMount(ctx app.Context) {
	url := app.Window().URL()
	c.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
	c.handleNextExercise()
	c.bindKeys()
	c.bindSwipes()
}

// The Render method is where the component appearance is defined.
func (c *Home) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.H2().Body(
			app.Text("What is the capital of "),
			app.If(c.question != "",
				app.Text(c.question),
			).Else(
				app.Text(""),
			),
			app.Text("?"),
		),
		app.H2().Body(
			app.If(c.answer != "",
				app.Text(c.answer),
			).Else(
				app.Text(""),
			),
		).Hidden(!c.isAnswerVisible),
		app.P().Body(
			app.Text("Good answers: "),
			app.Text(c.goodAnswers),
		),
		app.P().Body(
			app.Text("Bad answers: "),
			app.Text(c.badAnswers),
		),
		app.P().Body(
			app.Button().
				Text("⇧ View answer").
				OnClick(func(ctx app.Context, e app.Event) {
					c.handleViewAnswer()
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
			app.Button().
				Text("Next exercise ⇩").
				OnClick(func(ctx app.Context, e app.Event) {
					// only allow if next exercise was preloaded (to avoid double clicks)
					if c.isNextPreloaded == true {
						c.handleNextExercise()
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
					if c.isNextPreloaded == true {
						c.handleBadAnswer()
					}
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
			app.Button().
				Text("Good answer ⇨").
				OnClick(func(ctx app.Context, e app.Event) {
					// only allow if next exercise was preloaded (to avoid double clicks)
					if c.isNextPreloaded == true {
						c.handleGoodAnswer()
					}
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
		),
	)
}

func (c *Home) bindKeys() {
	app.Window().AddEventListener("keyup", func(ctx app.Context, e app.Event) {
		// bind actions to keyboard shortcuts
		switch e.Get("code").String() {
		case "Space":
			if c.isAnswerVisible == true {
				// only allow if next exercise was preloaded (to avoid double clicks)
				if c.isNextPreloaded == true {
					c.handleNextExercise()
				}
			} else {
				c.handleViewAnswer()
			}
		case "KeyV", "ArrowUp":
			c.handleViewAnswer()
		case "KeyG", "ArrowRight":
			// only allow if next exercise was preloaded (to avoid double clicks)
			if c.isNextPreloaded == true {
				c.handleGoodAnswer()
			}
		case "KeyB", "ArrowLeft":
			// only allow if next exercise was preloaded (to avoid double clicks)
			if c.isNextPreloaded == true {
				c.handleBadAnswer()
			}
		case "KeyN", "ArrowDown":
			// only allow if next exercise was preloaded (to avoid double clicks)
			if c.isNextPreloaded == true {
				c.handleNextExercise()
			}
		}
	})
}

func (c *Home) bindSwipes() {
	app.Window().AddEventListener("swiped-left", func(ctx app.Context, e app.Event) {
		// only allow if next exercise was preloaded (to avoid double clicks)
		if c.isNextPreloaded == true {
			c.handleBadAnswer()
		}
	})
	app.Window().AddEventListener("swiped-right", func(ctx app.Context, e app.Event) {
		// only allow if next exercise was preloaded (to avoid double clicks)
		if c.isNextPreloaded == true {
			c.handleGoodAnswer()
		}
	})
	app.Window().AddEventListener("swiped-up", func(ctx app.Context, e app.Event) {
		// only allow if next exercise was preloaded (to avoid double clicks)
		if c.isNextPreloaded == true {
			c.handleNextExercise()
		}
	})
	app.Window().AddEventListener("swiped-down", func(ctx app.Context, e app.Event) {
		c.handleViewAnswer()
	})
}

func (c *Home) handleNextExercise() {
	c.isAnswerVisible = false

	if c.isNextPreloaded == false {
		exercise := c.fetchNext()
		c.exerciseId = exercise.Id
		c.question = exercise.Question
		c.answer = exercise.Answer
		c.goodAnswers = exercise.GoodAnswers
		c.badAnswers = exercise.BadAnswers
		app.Log("displayed initial exercise")
	} else {
		c.isNextPreloaded = false
		c.exerciseId = c.nextExerciseId
		c.question = c.nextQuestion
		c.answer = c.nextAnswer
		c.goodAnswers = c.nextGoodAnswers
		c.badAnswers = c.nextBadAnswers
		app.Log("displayed preloaded exercise")
	}

	go func() {
		exercise := c.fetchNext()
		c.nextExerciseId = exercise.Id
		c.nextQuestion = exercise.Question
		c.nextAnswer = exercise.Answer
		c.nextGoodAnswers = exercise.GoodAnswers
		c.nextBadAnswers = exercise.BadAnswers
		c.isNextPreloaded = true
		app.Log("preloaded")
	}()
}

func (c *Home) fetchNext() models.Exercise {
	exercise, err := c.api.FetchNextExercise()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch next exercise: %w", err))
	}

	// dummy way of avoiding duplicates todo move to the API
	if exercise.Id == c.exerciseId {
		return c.fetchNext()
	}

	return exercise
}

func (c *Home) handleViewAnswer() {
	c.isAnswerVisible = true
}

func (c *Home) handleGoodAnswer() {
	exercise := models.Exercise{Id: c.exerciseId}
	go func() {
		// increment in the background
		if err := c.api.IncrementGoodAnswers(exercise); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	c.handleNextExercise()
}

func (c *Home) handleBadAnswer() {
	exercise := models.Exercise{Id: c.exerciseId}
	go func() {
		// increment in the background
		if err := c.api.IncrementBadAnswers(exercise); err != nil {
			app.Log(fmt.Errorf("failed to increment bad answers: %w", err))
		}
	}()
	c.handleNextExercise()
}
