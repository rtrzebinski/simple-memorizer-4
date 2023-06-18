package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal"
	"github.com/rtrzebinski/simple-memorizer-4/internal/http/rest"
	"github.com/rtrzebinski/simple-memorizer-4/internal/memorizer"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"net/http"
	"net/url"
	"strconv"
)

const PathLearn = "/learn"

// A Learn component
type Learn struct {
	app.Compo
	s         *internal.Service
	memorizer memorizer.Service

	// component vars
	lesson          models.Lesson
	exercise        models.Exercise
	isAnswerVisible bool

	// Window events unsubscribers, to be called on component dismount
	// this is needed so Window events are not piled on each component mounting
	unsubscribers []func()
}

// The OnMount method is run once component is mounted
func (c *Learn) OnMount(ctx app.Context) {
	u := app.Window().URL()

	// create a service, because if go-app lib limitations it can not be injected from main
	r := rest.NewReader(rest.NewClient(&http.Client{}, u.Host, u.Scheme))
	w := rest.NewWriter(rest.NewClient(&http.Client{}, u.Host, u.Scheme))
	c.s = internal.NewService(r, w)

	lessonId, err := strconv.Atoi(u.Query().Get("lesson_id"))
	if err != nil {
		app.Log("invalid lesson_id")

		return
	}

	c.lesson = models.Lesson{Id: lessonId}
	c.hydrateLesson()

	exercisesOfLesson, err := c.s.FetchExercisesOfLesson(c.lesson)
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch exercises of lesson: %w", err))

		return
	}

	// convert into map that is needed for a memorizer service
	var exercises = make(map[int]models.Exercise)
	for _, e := range exercisesOfLesson {
		exercises[e.Id] = e
	}

	c.memorizer.Init(exercises)

	c.handleNextExercise()
	c.bindKeys()
	c.bindSwipes()
}

// HydrateLesson in go routine
func (c *Learn) hydrateLesson() {
	go func() {
		err := c.s.HydrateLesson(&c.lesson)
		if err != nil {
			app.Log(fmt.Errorf("failed to hydrate lesson: %w", err))
		}
		// needs to be run manually so UI reflects the change
		c.Update()
	}()
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
			app.If(c.exercise.Question != "",
				app.Text(c.exercise.Question),
			).Else(
				app.Text(""),
			),
		),
		app.P().Body(
			app.Text("Answer: "),
			app.If(c.exercise.Answer != "",
				app.Text(c.exercise.Answer),
			).Else(
				app.Text(""),
			),
		).Hidden(!c.isAnswerVisible),
		app.P().Body(
			app.Text("Bad answers: "),
			app.Text(c.exercise.AnswersProjection.BadAnswers),
		),
		app.P().Body(
			app.Text("Good answers: "),
			app.Text(c.exercise.AnswersProjection.GoodAnswers),
		),
		app.P().Body(
			app.Text("Good answers %: "),
			app.Text(c.exercise.AnswersProjection.GoodAnswersPercent()),
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
					c.handleNextExercise()
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
		),
		app.P().Body(
			app.Button().
				Text("⇦ Bad answer").
				OnClick(func(ctx app.Context, e app.Event) {
					c.handleBadAnswer()
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
			app.Button().
				Text("Good answer ⇨").
				OnClick(func(ctx app.Context, e app.Event) {
					c.handleGoodAnswer()
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
				c.handleNextExercise()
			} else {
				c.handleViewAnswer()
			}
		case "KeyV", "ArrowUp":
			c.handleViewAnswer()
		case "KeyG", "ArrowRight":
			c.handleGoodAnswer()
		case "KeyB", "ArrowLeft":
			c.handleBadAnswer()
		case "KeyN", "ArrowDown":
			c.handleNextExercise()
		}
	})

	c.unsubscribers = append(c.unsubscribers, f)
}

func (c *Learn) bindSwipes() {
	var f func()

	f = app.Window().AddEventListener("swiped-left", func(ctx app.Context, e app.Event) {
		c.handleBadAnswer()
	})
	c.unsubscribers = append(c.unsubscribers, f)

	f = app.Window().AddEventListener("swiped-right", func(ctx app.Context, e app.Event) {
		c.handleGoodAnswer()
	})
	c.unsubscribers = append(c.unsubscribers, f)

	f = app.Window().AddEventListener("swiped-up", func(ctx app.Context, e app.Event) {
		c.handleNextExercise()
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
	c.exercise = c.memorizer.Next(c.exercise)
}

func (c *Learn) handleViewAnswer() {
	c.isAnswerVisible = true
}

func (c *Learn) handleGoodAnswer() {
	// copy so go routine will not rely on dynamic c.exercise
	exCopy := c.exercise
	// store answer in the background
	go func() {
		if err := c.s.StoreAnswer(&models.Answer{
			Exercise: &exCopy,
			Type:     models.Good,
		}); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	// exercise will be passed back to memorizer, it needs the correct answers count
	c.exercise.AnswersProjection.GoodAnswers++
	c.handleNextExercise()
}

func (c *Learn) handleBadAnswer() {
	// copy so go routine will not rely on dynamic c.exercise
	exCopy := c.exercise
	// store answer in the background
	go func() {
		if err := c.s.StoreAnswer(&models.Answer{
			Exercise: &exCopy,
			Type:     models.Bad,
		}); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	// exercise will be passed back to memorizer, it needs the correct answers count
	c.exercise.AnswersProjection.BadAnswers++
	c.handleNextExercise()
}
