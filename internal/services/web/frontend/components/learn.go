package components

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/memorizer"
)

const PathLearn = "/learn"

// Learn is a component that displays learning page
type Learn struct {
	app.Compo
	c               APIClient
	memorizer       memorizer.Memorizer
	lesson          frontend.Lesson
	exercise        frontend.Exercise
	isAnswerVisible bool
	user            *frontend.User
	edit            *ExerciseEdit // edit component for exercises
	editMode        bool          // used to disable answering while editing exercise

	// Window events unsubscribers, to be called on component dismount
	// this is needed so Window events are not piled on each component mounting
	unsubscribers []func()
}

// NewLearn creates a new Learn component
func NewLearn(c APIClient) *Learn {
	learn := &Learn{
		c:        c,
		user:     &frontend.User{},
		editMode: false,
	}

	learn.edit = NewExerciseEdit(c, learn)

	return learn
}

// The OnMount method is run once component is mounted
func (compo *Learn) OnMount(ctx app.Context) {
	// auth check
	user, err := auth.User(ctx)
	if err != nil {
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}
	compo.user = user

	u := app.Window().URL()

	lessonId, err := strconv.Atoi(u.Query().Get("lesson_id"))
	if err != nil {
		app.Log("invalid lesson_id")

		return
	}

	compo.lesson = frontend.Lesson{Id: lessonId}
	compo.hydrateLesson(ctx)

	oldestExerciseID, err := strconv.Atoi(u.Query().Get("oldest_exercise_id"))
	if err != nil {
		oldestExerciseID = 1
	}

	accessToken, err := auth.Token(ctx)
	if err != nil {
		slog.Error("failed to get token", "err", err)
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}
	exercisesOfLesson, err := compo.c.FetchExercises(ctx, compo.lesson, oldestExerciseID, accessToken)
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch exercises of lesson: %w", err))

		return
	}

	// convert into map that is needed for a memorizer service
	var exercises = make(map[int]frontend.Exercise)
	for _, e := range exercisesOfLesson {
		exercises[e.Id] = e
	}

	compo.memorizer.Init(exercises)

	compo.handleNextExercise()
	compo.bindKeys(ctx)
	compo.bindSwipes(ctx)
}

// HydrateLesson in go routine
func (compo *Learn) hydrateLesson(ctx app.Context) {
	go func() {
		accessToken, err := auth.Token(ctx)
		if err != nil {
			slog.Error("failed to get token", "err", err)
			ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
		}
		err = compo.c.HydrateLesson(ctx, &compo.lesson, accessToken)
		if err != nil {
			app.Log(fmt.Errorf("failed to hydrate lesson: %w", err))
		}
	}()
}

// The OnDismount method is run once component is dismounted.
func (compo *Learn) OnDismount() {
	compo.runUnsubscribeFunctions()
}

// runUnsubscribeFunctions unsubscribe from registered Window events
func (compo *Learn) runUnsubscribeFunctions() {
	for _, f := range compo.unsubscribers {
		f()
	}
}

// The Render method is where the component appearance is defined.
func (compo *Learn) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.P().Body(
				app.A().Href(PathHome).Text("Home"),
				app.Text(" | "),
				app.A().Href(PathLessons).Text("Lessons"),
				app.Text(" | "),
				app.A().Href(PathAuthLogout).Text("Logout"),
				app.Text(" | "),
				app.Text(app.Getenv("version")),
			),
		),
		app.P().Body(
			app.Button().Text("Show exercises").OnClick(compo.handleShowExercises),
			app.Button().Text("Edit exercise").OnClick(compo.handleEditExercise),
		),
		app.Div().Body(compo.edit.Render()),
		app.P().Body(
			app.Text("Question: "),
			app.If(compo.exercise.Question != "", func() app.UI {
				return app.Text(compo.exercise.Question)
			}).Else(func() app.UI {
				return app.Text("")
			}),
		).Hidden(compo.editMode),
		app.P().Body(
			app.Text("Answer: "),
			app.If(compo.exercise.Answer != "", func() app.UI {
				return app.Text(compo.exercise.Answer)
			}).Else(func() app.UI {
				return app.Text("")
			}),
		).Hidden(!compo.isAnswerVisible || compo.editMode),
		app.P().Body(
			app.Button().
				Text("⇧ View answer").
				OnClick(func(ctx app.Context, e app.Event) {
					compo.handleViewAnswer()
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
			app.Button().
				Text("Next exercise ⇩").
				OnClick(func(ctx app.Context, e app.Event) {
					compo.handleNextExercise()
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
		).Hidden(compo.editMode),
		app.P().Body(
			app.Button().
				Text("⇦ Bad answer").
				OnClick(func(ctx app.Context, e app.Event) {
					compo.handleBadAnswer(ctx)
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
			app.Button().
				Text("Good answer ⇨").
				OnClick(func(ctx app.Context, e app.Event) {
					compo.handleGoodAnswer(ctx)
				}).
				Style("margin-right", "10px").
				Style("font-size", "15px"),
		).Hidden(compo.editMode),
	)
}

// bindKeys binds keyboard events to actions
func (compo *Learn) bindKeys(ctx app.Context) {
	hFn := func(this app.Value, args []app.Value) any {
		event := args[0]
		app.Log("key pressed: " + event.Get("code").String())
		ctx.Dispatch(func(ctx app.Context) {
			// bind actions to keyboard shortcuts
			switch event.Get("code").String() {
			case "Space":
				if compo.isAnswerVisible == true {
					compo.handleNextExercise()
				} else {
					compo.handleViewAnswer()
				}
			case "KeyV", "ArrowUp":
				compo.handleViewAnswer()
			case "KeyG", "ArrowRight":
				compo.handleGoodAnswer(ctx)
			case "KeyB", "ArrowLeft":
				compo.handleBadAnswer(ctx)
			case "KeyN", "ArrowDown":
				compo.handleNextExercise()
			}
		})

		return nil
	}

	fo := app.FuncOf(hFn)

	app.Window().Call("addEventListener", "keyup", fo)

	compo.unsubscribers = append(compo.unsubscribers, func() {
		app.Window().Call("removeEventListener", "keyup", fo)
		fo.Release()
	})
}

// bindSwipes binds swipe events to actions
func (compo *Learn) bindSwipes(ctx app.Context) {
	compo.bindSwipe(ctx, "swiped-left", func(ctx app.Context) {
		compo.handleBadAnswer(ctx)
	})
	compo.bindSwipe(ctx, "swiped-right", func(ctx app.Context) {
		compo.handleGoodAnswer(ctx)
	})
	compo.bindSwipe(ctx, "swiped-up", func(ctx app.Context) {
		compo.handleNextExercise()
	})
	compo.bindSwipe(ctx, "swiped-down", func(ctx app.Context) {
		compo.handleViewAnswer()
	})
}

// bindSwipe binds swipe event to an action
func (compo *Learn) bindSwipe(ctx app.Context, eventName string, v func(app.Context)) {
	hFn := func(this app.Value, args []app.Value) any {
		ctx.Dispatch(func(ctx app.Context) {
			v(ctx)
		})

		return nil
	}

	fo := app.FuncOf(hFn)

	app.Window().Call("addEventListener", eventName, fo)

	compo.unsubscribers = append(compo.unsubscribers, func() {
		app.Window().Call("removeEventListener", eventName, fo)
		fo.Release()
	})
}

// handleShowExercises start learning a current lesson
func (compo *Learn) handleShowExercises(ctx app.Context, _ app.Event) {
	app.Log("handleShowExercises")
	u, _ := url.Parse(PathExercises)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(compo.lesson.Id))
	u.RawQuery = params.Encode()

	ctx.NavigateTo(u)
}

// handleNextExercise moves to the next exercise
func (compo *Learn) handleNextExercise() {
	app.Log("handleNextExercise")
	compo.isAnswerVisible = false
	compo.exercise = compo.memorizer.Next(compo.exercise)
}

// handleViewAnswer shows the answer
func (compo *Learn) handleViewAnswer() {
	app.Log("handleViewAnswer")
	compo.isAnswerVisible = true
}

// handleGoodAnswer increments good answers and moves to the next exercise
func (compo *Learn) handleGoodAnswer(ctx app.Context) {
	app.Log("handleGoodAnswer")
	// copy so go routine will not rely on dynamic compo.exercise
	exCopy := compo.exercise
	// save answer in the background
	go func() {
		accessToken, err := auth.Token(ctx)
		if err != nil {
			slog.Error("failed to get token", "err", err)
			ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
		}
		if err := compo.c.StoreResult(ctx, frontend.Result{
			Exercise: &exCopy,
			Type:     frontend.Good,
		}, accessToken); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	// exercise will be passed back to memorizer, needs the updated projection
	compo.exercise.RegisterGoodAnswer()
	compo.handleNextExercise()
}

// handleBadAnswer increments bad answers and moves to the next exercise
func (compo *Learn) handleBadAnswer(ctx app.Context) {
	app.Log("handleBadAnswer")
	// copy so go routine will not rely on dynamic compo.exercise
	exCopy := compo.exercise
	// save answer in the background
	go func() {
		accessToken, err := auth.Token(ctx)
		if err != nil {
			slog.Error("failed to get token", "err", err)
			ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
		}
		if err := compo.c.StoreResult(ctx, frontend.Result{
			Exercise: &exCopy,
			Type:     frontend.Bad,
		}, accessToken); err != nil {
			app.Log(fmt.Errorf("failed to increment good answers: %w", err))
		}
	}()
	// exercise will be passed back to memorizer, needs the updated projection
	compo.exercise.RegisterBadAnswer()
	compo.handleNextExercise()
}

// handleEditExercise is called when the "edit exercise" button is clicked
func (compo *Learn) handleEditExercise(_ app.Context, _ app.Event) {
	compo.edit.edit(compo.exercise) // set edit component to edit mode
	compo.runUnsubscribeFunctions() // remove all previous event listeners (unbind keys and swipes)
	compo.editMode = true           // disable answering while editing
}

// exerciseEditDone is called when the edit form is submitted or cancelled
func (compo *Learn) exerciseEditDone(ctx app.Context) {
	compo.exercise.Question = compo.edit.inputQuestion // update question
	compo.exercise.Answer = compo.edit.inputAnswer     // update answer
	compo.bindKeys(ctx)                                // rebind keys after edit is done
	compo.bindSwipes(ctx)                              // rebind swipes after edit is done
	compo.editMode = false                             // enable answering after edit is done
}

// getLesson returns the current lesson
func (compo *Learn) getLesson() *frontend.Lesson {
	return &compo.lesson
}

// getExercises returns a slice of exercises from the memorizer
func (compo *Learn) getExercises() []*frontend.Exercise {
	return compo.memorizer.GetExercises()
}
