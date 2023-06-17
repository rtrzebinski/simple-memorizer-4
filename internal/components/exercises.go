package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal"
	"github.com/rtrzebinski/simple-memorizer-4/internal/http/rest"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validation"
	"net/http"
	"net/url"
	"strconv"
)

const PathExercises = "/exercises"

type Exercises struct {
	app.Compo
	s *internal.Service

	// component vars
	lesson models.Lesson
	rows   []*ExerciseRow

	// store exercise form
	formVisible         bool
	validationErrors    []error
	inputId             int
	inputQuestion       string
	inputAnswer         string
	storeButtonDisabled bool
}

// The OnMount method is run once component is mounted
func (c *Exercises) OnMount(ctx app.Context) {
	u := app.Window().URL()

	// create a service, because if go-app lib limitations it can not be injected from main
	r := rest.NewReader(rest.NewClient(&http.Client{}, u.Host, u.Scheme))
	w := rest.NewWriter(rest.NewClient(&http.Client{}, u.Host, u.Scheme))
	c.s = internal.NewService(r, w)

	lessonId, err := strconv.Atoi(u.Query().Get("lesson_id"))
	if err != nil {
		app.Log(fmt.Errorf("failed to convert lesson_id: %w", err))
		return
	}

	c.lesson = models.Lesson{Id: lessonId}

	c.hydrateLesson()
	c.displayExercisesOfLesson()
}

// The Render method is where the component appearance is defined.
func (c *Exercises) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Button().Text("Start learning").OnClick(c.handleStartLearning).Disabled(c.lesson.ExerciseCount < 2),
			app.Button().Text("Add a new exercise").OnClick(c.handleAddExercise).Hidden(c.formVisible),
		),
		app.P().Body(
			app.Text("Lesson name: "),
			app.Text(c.lesson.Name),
		),
		app.P().Body(
			app.Text("Lesson description: "),
			app.Text(c.lesson.Description),
		),
		app.P().Body(
			app.Text("Exercises: "),
			app.Text(c.lesson.ExerciseCount),
		),
		app.Div().Body(
			app.H3().Text("Add a new exercise"),
			app.P().Body(
				app.Range(c.validationErrors).Slice(func(i int) app.UI {
					return app.Div().Body(
						app.Text(c.validationErrors[i].Error()),
						app.Br(),
					).Style("color", "red")
				}),
			),
			app.Textarea().ID("input_question").Cols(30).Rows(3).Placeholder("Question").
				Required(true).OnInput(c.ValueTo(&c.inputQuestion)).Text(c.inputQuestion),
			app.Br(),
			app.Textarea().ID("input_answer").Cols(30).Rows(3).Placeholder("Answer").
				Required(true).OnInput(c.ValueTo(&c.inputAnswer)).Text(c.inputAnswer),
			app.Br(),
			app.Button().Text("Store").OnClick(c.handleStore).Disabled(c.storeButtonDisabled),
			app.Button().Text("Cancel").OnClick(c.handleCancel),
		).Hidden(!c.formVisible),
		app.Div().Body(
			app.H3().Text("Exercises"),
			app.Table().Style("border", "1px solid black").Body(
				&ExerciseHeader{},
				app.Range(c.rows).Slice(func(i int) app.UI {
					if c.rows[i] != nil {
						return c.rows[i].Render()
					}
					return nil
				}),
			),
		),
	)
}

// handleStartLearning start learning a current lesson
func (c *Exercises) handleStartLearning(ctx app.Context, e app.Event) {
	u, _ := url.Parse(PathLearn)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(c.lesson.Id))
	u.RawQuery = params.Encode()

	ctx.NavigateTo(u)
}

// handleAddExercise display add exercise form
func (c *Exercises) handleAddExercise(ctx app.Context, e app.Event) {
	c.formVisible = true
}

// handleStore create new or update existing exercise
func (c *Exercises) handleStore(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// exercise to be stored
	exercise := models.Exercise{
		Id:       c.inputId,
		Question: c.inputQuestion,
		Answer:   c.inputAnswer,
		Lesson: &models.Lesson{
			Id: c.lesson.Id,
		},
	}

	// extract other questions to validate against
	var questions []string
	for i, row := range c.rows {
		if row != nil && c.rows[i].exercise.Id != c.inputId {
			questions = append(questions, c.rows[i].exercise.Question)
		}
	}

	// validate input
	validator := validation.ValidateStoreExercise(exercise, questions)
	if validator.Failed() {
		c.validationErrors = validator.Errors()

		return
	}

	// disable submit button to avoid duplicated requests
	c.storeButtonDisabled = true

	// store exercise
	err := c.s.StoreExercise(&exercise)
	if err != nil {
		app.Log(fmt.Errorf("failed to store exercise: %w", err))
	}

	// hide the input form on row edit, but keep open on adding new
	// because it is common to add a few exercises one after another
	if c.inputId != 0 {
		c.formVisible = false
	}

	// reset form
	c.resetForm()

	// refresh exercises list
	c.displayExercisesOfLesson()

	// refresh lesson details (exercises counter)
	c.hydrateLesson()
}

func (c *Exercises) handleCancel(ctx app.Context, e app.Event) {
	c.resetForm()
	c.formVisible = false
}

func (c *Exercises) resetForm() {
	c.inputId = 0
	c.inputQuestion = ""
	c.inputAnswer = ""
	c.validationErrors = nil
	c.storeButtonDisabled = false
}

func (c *Exercises) displayExercisesOfLesson() {
	exercises, err := c.s.FetchExercisesOfLesson(c.lesson)
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch exercises of lesson: %w", err))
	}

	// no entries in the database
	if len(exercises) == 0 {
		return
	}

	// find maxId so we know the rows slice capacity
	maxId := exercises[0].Id
	for _, row := range exercises {
		if row.Id > maxId {
			maxId = row.Id
		}
	}

	// add +1 to len as IDs from the DB are 1 indexed, while slices are 0 indexed,
	// so we need to shift by one to have space for the latest row
	c.rows = make([]*ExerciseRow, maxId+1)

	for _, row := range exercises {
		c.rows[row.Id] = &ExerciseRow{exercise: row, parent: c}
	}
}

func (c *Exercises) hydrateLesson() {
	err := c.s.HydrateLesson(&c.lesson)
	if err != nil {
		app.Log(fmt.Errorf("failed to hydrate lesson: %w", err))
		return
	}
}
