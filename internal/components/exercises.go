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

const PathExercises = "/exercises"

type Exercises struct {
	app.Compo
	reader *rest.Reader
	writer *rest.Writer
	lesson models.Lesson
	rows   []*ExerciseRow

	// store exercise form
	formVisible         bool
	validationError     string
	inputId             int
	inputQuestion       string
	inputAnswer         string
	storeButtonDisabled bool
}

// The OnMount method is run once component is mounted
func (c *Exercises) OnMount(ctx app.Context) {
	u := app.Window().URL()

	c.reader = rest.NewReader(rest.NewClient(&http.Client{}, u.Host, u.Scheme))
	c.writer = rest.NewWriter(rest.NewClient(&http.Client{}, u.Host, u.Scheme))

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
			app.Text("Lesson: "),
			app.Text(c.lesson.Name),
		),
		app.P().Body(
			app.Text("Exercises: "),
			app.Text(c.lesson.ExerciseCount),
		),
		app.Div().Body(
			app.H3().Text("Add a new exercise"),
			app.P().Text(c.validationError).Hidden(c.validationError == "").Style("color", "red"),
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

	// init empty validation errors
	c.validationError = ""

	// validate input - exercise question required
	if c.inputQuestion == "" {
		c.validationError = "Exercise question is required"
		return
	}

	// validate input - exercise answer required
	if c.inputAnswer == "" {
		c.validationError = "Exercise answer is required"
		return
	}

	// validate input - exercise question unique (for given lesson)
	for i, row := range c.rows {
		if row != nil && c.rows[i].exercise.Question == c.inputQuestion && c.rows[i].exercise.Id != c.inputId {
			c.validationError = "Exercise question must be unique"
			return
		}
	}

	// check before resetting the form
	isCreatingNew := c.inputId == 0

	// disable submit button to avoid duplicated requests
	c.storeButtonDisabled = true

	// store exercise
	err := c.writer.StoreExercise(models.Exercise{
		Id:       c.inputId,
		Question: c.inputQuestion,
		Answer:   c.inputAnswer,
		Lesson: &models.Lesson{
			Id: c.lesson.Id,
		},
	})
	if err != nil {
		app.Log(fmt.Errorf("failed to store exercise: %w", err))
	}

	// reset form
	c.resetForm()

	// hide the input form on row edit, but keep open on adding new
	// because it is common to add a few exercises one after another
	if isCreatingNew == false {
		c.formVisible = false
	}

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
	c.validationError = ""
	c.storeButtonDisabled = false
}

func (c *Exercises) displayExercisesOfLesson() {
	exercises, err := c.reader.FetchExercisesOfLesson(c.lesson)
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
	err := c.reader.HydrateLesson(&c.lesson)
	if err != nil {
		app.Log(fmt.Errorf("failed to hydrate lesson: %w", err))
		return
	}
}
