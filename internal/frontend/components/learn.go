package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/api"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/memorizer"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"net/url"
	"strconv"
)

const PathLearn = "/learn"

// Learn is a component that displays learning page
type Learn struct {
	app.Compo
	s         *api.Service
	memorizer memorizer.Service

	// component vars
	lesson          models.Lesson
	exercise        models.Exercise
	isAnswerVisible bool

	// Window events unsubscribers, to be called on component dismount
	// this is needed so Window events are not piled on each component mounting
	unsubscribers []func()
}

// NewLearn creates a new Learn component
func NewLearn(s *api.Service) *Learn {
	return &Learn{
		s: s,
	}
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
	c.hydrateLesson()

	exercisesOfLesson, err := c.s.FetchExercises(c.lesson)
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
	c.bindKeys(ctx)
	c.bindSwipes(ctx)
}

// HydrateLesson in go routine
func (c *Learn) hydrateLesson() {
	go func() {
		err := c.s.HydrateLesson(&c.lesson)
		if err != nil {
			app.Log(fmt.Errorf("failed to hydrate lesson: %w", err))
		}
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
			app.If(c.exercise.Question != "", func() app.UI {
				return app.Text(c.exercise.Question)
			}).Else(func() app.UI {
				return app.Text("")
			}),
		),
		app.P().Body(
			app.Text("Answer: "),
			app.If(c.exercise.Answer != "", func() app.UI {
				return app.Text(c.exercise.Answer)
			}).Else(func() app.UI {
				return app.Text("")
			}),
		).Hidden(!c.isAnswerVisible),
		app.P().Body(
			app.Text(c.exercise.ResultsProjection.BadAnswers),
			app.Text(" bad answers"),
			app.If(c.exercise.ResultsProjection.BadAnswers > 0 && c.exercise.ResultsProjection.BadAnswersToday > 0, func() app.UI {
				return app.Text(" (today: " + strconv.Itoa(c.exercise.ResultsProjection.BadAnswersToday) + ")")
			}).Else(func() app.UI {
				return app.Text(" (latest: " + c.exercise.ResultsProjection.LatestBadAnswer.Format("2 Jan 2006 15:04") + ")")
			}),
		),
		app.P().Body(
			app.Text(c.exercise.ResultsProjection.GoodAnswers),
			app.Text(" good answers"),

			app.If(c.exercise.ResultsProjection.GoodAnswers > 0 && c.exercise.ResultsProjection.GoodAnswersToday > 0, func() app.UI {
				return app.Text(" (today: " + strconv.Itoa(c.exercise.ResultsProjection.GoodAnswersToday) + ")")
			}).Else(func() app.UI {
				return app.Text(" (latest: " + c.exercise.ResultsProjection.LatestGoodAnswer.Format("2 Jan 2006 15:04") + ")")
			}),
		),
		app.P().Body(
			app.Text(c.exercise.ResultsProjection.GoodAnswersPercent()),
			app.Text("% of good answers"),
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

// bindKeys binds keyboard events to actions
func (c *Learn) bindKeys(ctx app.Context) {
	hFn := func(this app.Value, args []app.Value) any {
		event := args[0]
		app.Log("key pressed: " + event.Get("code").String())
		ctx.Dispatch(func(ctx app.Context) {
			// bind actions to keyboard shortcuts
			switch event.Get("code").String() {
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

		return nil
	}

	fo := app.FuncOf(hFn)

	app.Window().Call("addEventListener", "keyup", fo)

	c.unsubscribers = append(c.unsubscribers, func() {
		app.Window().Call("removeEventListener", "keyup", fo)
		fo.Release()
	})
}

// bindSwipes binds swipe events to actions
func (c *Learn) bindSwipes(ctx app.Context) {
	c.bindSwipe(ctx, "swiped-left", func(ctx app.Context) {
		c.handleBadAnswer()
	})
	c.bindSwipe(ctx, "swiped-right", func(ctx app.Context) {
		c.handleGoodAnswer()
	})
	c.bindSwipe(ctx, "swiped-up", func(ctx app.Context) {
		c.handleNextExercise()
	})
	c.bindSwipe(ctx, "swiped-down", func(ctx app.Context) {
		c.handleViewAnswer()
	})
}

// bindSwipe binds swipe event to an action
func (c *Learn) bindSwipe(ctx app.Context, eventName string, v func(app.Context)) {
	hFn := func(this app.Value, args []app.Value) any {
		ctx.Dispatch(func(ctx app.Context) {
			v(ctx)
		})

		return nil
	}

	fo := app.FuncOf(hFn)

	app.Window().Call("addEventListener", eventName, fo)

	c.unsubscribers = append(c.unsubscribers, func() {
		app.Window().Call("removeEventListener", eventName, fo)
		fo.Release()
	})
}

// handleShowExercises start learning a current lesson
func (c *Learn) handleShowExercises(ctx app.Context, _ app.Event) {
	app.Log("handleShowExercises")
	u, _ := url.Parse(PathExercises)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(c.lesson.Id))
	u.RawQuery = params.Encode()

	ctx.NavigateTo(u)
}

// handleNextExercise moves to the next exercise
func (c *Learn) handleNextExercise() {
	app.Log("handleNextExercise")
	c.isAnswerVisible = false
	c.exercise = c.memorizer.Next(c.exercise)
}

// handleViewAnswer shows the answer
func (c *Learn) handleViewAnswer() {
	app.Log("handleViewAnswer")
	c.isAnswerVisible = true
}

// handleGoodAnswer increments good answers and moves to the next exercise
func (c *Learn) handleGoodAnswer() {
	app.Log("handleGoodAnswer")
	// copy so go routine will not rely on dynamic c.exercise
	exCopy := c.exercise
	// save answer in the background
	go func() {
		if err := c.s.StoreResult(&models.Result{
			Exercise: &exCopy,
			Type:     models.Good,
		}); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	// exercise will be passed back to memorizer, needs the updated projection
	c.exercise.ResultsProjection.RegisterGoodAnswer()
	c.handleNextExercise()
}

// handleBadAnswer increments bad answers and moves to the next exercise
func (c *Learn) handleBadAnswer() {
	app.Log("handleBadAnswer")
	// copy so go routine will not rely on dynamic c.exercise
	exCopy := c.exercise
	// save answer in the background
	go func() {
		if err := c.s.StoreResult(&models.Result{
			Exercise: &exCopy,
			Type:     models.Bad,
		}); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	// exercise will be passed back to memorizer, needs the updated projection
	c.exercise.ResultsProjection.RegisterBadAnswer()
	c.handleNextExercise()
}
