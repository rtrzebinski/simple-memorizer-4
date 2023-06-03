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

var pathExercises = "/exercises"

type Exercises struct {
	app.Compo
	api    *rest.Client
	lesson models.Lesson
	rows   []*ExerciseRow

	addExerciseFormVisible bool
	inputId                int
	inputQuestion          string
	inputAnswer            string
	storeButtonDisabled    bool
}

// The OnMount method is run once component is mounted
func (c *Exercises) OnMount(ctx app.Context) {
	url := app.Window().URL()

	lessonId, err := strconv.Atoi(url.Query().Get("lesson_id"))
	if err != nil {
		app.Log("invalid lesson_id")
		return
	}
	c.lesson = models.Lesson{Id: lessonId}

	c.api = rest.NewClient(&http.Client{}, url.Host, url.Scheme)
	c.displayExercisesOfLesson()
}

// The Render method is where the component appearance is defined.
func (c *Exercises) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Button().Text("Start learning").OnClick(c.handleStartLearning),
			app.Button().Text("Add a new exercise").OnClick(c.handleAddExercise).Hidden(c.addExerciseFormVisible),
		),
		app.Div().Body(
			app.H3().Text("Add a new exercise"),
			app.Textarea().ID("input_question").Cols(30).Rows(3).Placeholder("Question").
				Required(true).OnInput(c.ValueTo(&c.inputQuestion)).Text(c.inputQuestion),
			app.Br(),
			app.Textarea().ID("input_answer").Cols(30).Rows(3).Placeholder("Answer").
				Required(true).OnInput(c.ValueTo(&c.inputAnswer)).Text(c.inputAnswer),
			app.Br(),
			app.Button().Text("Store").OnClick(c.handleStore).Disabled(c.storeButtonDisabled),
			app.Button().Text("Cancel").OnClick(c.handleCancel),
		).Hidden(!c.addExerciseFormVisible),
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
	u, _ := url.Parse(pathLearn)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(c.lesson.Id))
	u.RawQuery = params.Encode()

	ctx.NavigateTo(u)
}

// handleAddExercise display add exercise form
func (c *Exercises) handleAddExercise(ctx app.Context, e app.Event) {
	c.addExerciseFormVisible = true
}

// handleStore create new or update existing exercise
func (c *Exercises) handleStore(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// todo implement input validation
	if c.inputQuestion == "" || c.inputAnswer == "" {
		app.Log("invalid input")
		return
	}

	isCreatingNew := c.inputId == 0

	c.storeButtonDisabled = true

	err := c.api.StoreExercise(models.Exercise{
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

	c.inputId = 0
	c.inputQuestion = ""
	c.inputAnswer = ""

	c.displayExercisesOfLesson()

	c.storeButtonDisabled = false

	// hide the input form on row edit, but keep open on adding new
	// because it is common to add a few exercises one after another
	if isCreatingNew == false {
		c.addExerciseFormVisible = false
	}
}

// handleCancel cleanup input UI
func (c *Exercises) handleCancel(ctx app.Context, e app.Event) {
	c.inputId = 0
	c.inputQuestion = ""
	c.inputAnswer = ""
	c.addExerciseFormVisible = false
}

func (c *Exercises) displayExercisesOfLesson() {
	exercises, err := c.api.FetchExercisesOfLesson(c.lesson)
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
