package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"net/http"
)

type Exercises struct {
	app.Compo
	api  *frontend.ApiClient
	rows []*ExerciseRow

	inputId             int
	inputQuestion       string
	inputAnswer         string
	storeButtonDisabled bool
}

// The OnMount method is run once component is mounted
func (c *Exercises) OnMount(ctx app.Context) {
	url := app.Window().URL()
	c.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
	c.displayAllExercises()
}

// The Render method is where the component appearance is defined.
func (c *Exercises) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.Div().Body(
			app.H2().Text("Store exercise"),
			app.Textarea().ID("input_question").Cols(30).Rows(3).Placeholder("Question").
				Required(true).OnInput(c.ValueTo(&c.inputQuestion)).Text(c.inputQuestion),
			app.Br(),
			app.Textarea().ID("input_answer").Cols(30).Rows(3).Placeholder("Answer").
				Required(true).OnInput(c.ValueTo(&c.inputAnswer)).Text(c.inputAnswer),
			app.Br(),
			app.Button().Text("Store").OnClick(c.handleStore).Disabled(c.storeButtonDisabled),
			app.Button().Text("Cancel").OnClick(c.handleCancel),
		),
		app.Div().Body(
			app.H2().Text("All exercises"),
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

// handleStore create new or update existing exercise
func (c *Exercises) handleStore(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// todo implement input validation
	if c.inputQuestion == "" || c.inputAnswer == "" {
		app.Log("invalid input")
		return
	}

	c.storeButtonDisabled = true

	err := c.api.StoreExercise(models.Exercise{
		Id:       c.inputId,
		Question: c.inputQuestion,
		Answer:   c.inputAnswer,
	})
	if err != nil {
		app.Log(fmt.Errorf("failed to store exercise: %w", err))
	}

	c.inputId = 0
	c.inputQuestion = ""
	c.inputAnswer = ""

	c.displayAllExercises()

	c.storeButtonDisabled = false
}

// handleCancel cleanup input UI
func (c *Exercises) handleCancel(ctx app.Context, e app.Event) {
	c.inputId = 0
	c.inputQuestion = ""
	c.inputAnswer = ""
}

func (c *Exercises) displayAllExercises() {
	exercises, err := c.api.FetchAllExercises()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch all exercises: %w", err))
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
