package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/http/rest"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"net/http"
	"net/url"
	"strconv"
)

const PathLearn = "/learn"

// A Learn component
type Learn struct {
	app.Compo
	reader *rest.Reader
	writer *rest.Writer
	lesson models.Lesson

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

	// Window events unsubscribers, to be called on component dismount
	// this is needed so Window events are not piled on each component mounting
	unsubscribers []func()
}

// The OnMount method is run once component is mounted
func (c *Learn) OnMount(ctx app.Context) {
	u := app.Window().URL()

	lessonId, err := strconv.Atoi(u.Query().Get("lesson_id"))
	if err != nil {
		app.Log("invalid lesson_id")
		return
	}
	c.lesson = models.Lesson{Id: lessonId}

	c.reader = rest.NewReader(rest.NewClient(&http.Client{}, u.Host, u.Scheme))
	c.writer = rest.NewWriter(rest.NewClient(&http.Client{}, u.Host, u.Scheme))
	err = c.reader.HydrateLesson(&c.lesson)
	if err != nil {
		app.Log(fmt.Errorf("failed to hydrate lesson: %w", err))
	}
	c.handleNextExercise()
	c.bindKeys()
	c.bindSwipes()
}

// The OnDismount method is run once component is dismounted.
func (c *Learn) OnDismount() {
	// unsubscribe from registered Window events
	for _, f := range c.unsubscribers {
		f()
	}
}

// The Render method is where the component appearance is defined.
func (c *Learn) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Button().Text("Show exercises").OnClick(c.handleShowExercises),
		),
		app.P().Body(
			app.Text("Lesson name: "+c.lesson.Name),
		),
		app.P().Body(
			app.Text("Lesson description: "+c.lesson.Description),
		),
		app.P().Body(
			app.Text("Question: "),
			app.If(c.question != "",
				app.Text(c.question),
			).Else(
				app.Text(""),
			),
		),
		app.P().Body(
			app.Text("Answer: "),
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

func (c *Learn) bindKeys() {
	var f func()

	f = app.Window().AddEventListener("keyup", func(ctx app.Context, e app.Event) {
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

	c.unsubscribers = append(c.unsubscribers, f)
}

func (c *Learn) bindSwipes() {
	var f func()

	f = app.Window().AddEventListener("swiped-left", func(ctx app.Context, e app.Event) {
		// only allow if next exercise was preloaded (to avoid double clicks)
		if c.isNextPreloaded == true {
			c.handleBadAnswer()
		}
	})
	c.unsubscribers = append(c.unsubscribers, f)

	f = app.Window().AddEventListener("swiped-right", func(ctx app.Context, e app.Event) {
		// only allow if next exercise was preloaded (to avoid double clicks)
		if c.isNextPreloaded == true {
			c.handleGoodAnswer()
		}
	})
	c.unsubscribers = append(c.unsubscribers, f)

	f = app.Window().AddEventListener("swiped-up", func(ctx app.Context, e app.Event) {
		// only allow if next exercise was preloaded (to avoid double clicks)
		if c.isNextPreloaded == true {
			c.handleNextExercise()
		}
	})
	c.unsubscribers = append(c.unsubscribers, f)

	f = app.Window().AddEventListener("swiped-down", func(ctx app.Context, e app.Event) {
		c.handleViewAnswer()
	})
	c.unsubscribers = append(c.unsubscribers, f)
}

// handleShowExercises start learning a current lesson
func (c *Learn) handleShowExercises(ctx app.Context, e app.Event) {
	u, _ := url.Parse(PathExercises)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(c.lesson.Id))
	u.RawQuery = params.Encode()

	ctx.NavigateTo(u)
}

func (c *Learn) handleNextExercise() {
	c.isAnswerVisible = false

	if c.isNextPreloaded == false {
		exercise := c.FetchNextExercise()
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
		exercise := c.FetchNextExercise()
		c.nextExerciseId = exercise.Id
		c.nextQuestion = exercise.Question
		c.nextAnswer = exercise.Answer
		c.nextGoodAnswers = exercise.GoodAnswers
		c.nextBadAnswers = exercise.BadAnswers
		c.isNextPreloaded = true
		app.Log("preloaded")
	}()
}

func (c *Learn) FetchNextExercise() models.Exercise {
	exercise, err := c.reader.FetchRandomExerciseOfLesson(c.lesson)
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch next exercise: %w", err))
	}

	// dummy way of avoiding duplicates todo move to the API
	if exercise.Id == c.exerciseId {
		return c.FetchNextExercise()
	}

	return exercise
}

func (c *Learn) handleViewAnswer() {
	c.isAnswerVisible = true
}

func (c *Learn) handleGoodAnswer() {
	exercise := models.Exercise{Id: c.exerciseId}
	go func() {
		// increment in the background
		if err := c.writer.IncrementGoodAnswers(exercise); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	c.handleNextExercise()
}

func (c *Learn) handleBadAnswer() {
	exercise := models.Exercise{Id: c.exerciseId}
	go func() {
		// increment in the background
		if err := c.writer.IncrementBadAnswers(exercise); err != nil {
			app.Log(fmt.Errorf("failed to increment bad answers: %w", err))
		}
	}()
	c.handleNextExercise()
}
