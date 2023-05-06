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
func (h *Exercises) OnMount(ctx app.Context) {
	url := app.Window().URL()
	h.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
	h.fetchAllExercises()
}

// The Render method is where the component appearance is defined.
func (h *Exercises) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.Div().Body(
			app.H2().Text("Store exercise"),
			app.Textarea().ID("input_question").Cols(30).Rows(3).Placeholder("Question").
				Required(true).OnInput(h.ValueTo(&h.inputQuestion)).Text(h.inputQuestion),
			app.Br(),
			app.Textarea().ID("input_answer").Cols(30).Rows(3).Placeholder("Answer").
				Required(true).OnInput(h.ValueTo(&h.inputAnswer)).Text(h.inputAnswer),
			app.Br(),
			app.Button().Text("Store").OnClick(h.handleStore).Disabled(h.storeButtonDisabled),
			app.Button().Text("Cancel").OnClick(h.handleCancel),
		),
		app.Div().Body(
			app.H2().Text("All exercises"),
			app.Table().Style("border", "1px solid black").Body(
				&ExerciseHeader{},
				app.Range(h.rows).Slice(func(i int) app.UI {
					if h.rows[i] != nil {
						return h.rows[i].Render()
					}
					return nil
				}),
			),
		),
	)
}

// handleStore create new or update existing exercise
func (h *Exercises) handleStore(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// todo implement input validation
	if h.inputQuestion == "" || h.inputAnswer == "" {
		app.Log("invalid input")
		return
	}

	h.storeButtonDisabled = true

	err := h.api.StoreExercise(models.Exercise{
		Id:       h.inputId,
		Question: h.inputQuestion,
		Answer:   h.inputAnswer,
	})
	if err != nil {
		app.Log(fmt.Errorf("failed to store exercise: %w", err))
	}

	h.inputId = 0
	h.inputQuestion = ""
	h.inputAnswer = ""

	h.fetchAllExercises()

	h.storeButtonDisabled = false
}

// handleCancel cleanup exercise input UI
func (h *Exercises) handleCancel(ctx app.Context, e app.Event) {
	h.inputId = 0
	h.inputQuestion = ""
	h.inputAnswer = ""
}

func (h *Exercises) fetchAllExercises() {
	exercises, err := h.api.FetchAllExercises()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch all exercises: %w", err))
	}

	// no exercises in the database
	if len(exercises) == 0 {
		return
	}

	// find maxId so we know the rows slice capacity
	maxId := exercises[0].Id
	for _, exercise := range exercises {
		if exercise.Id > maxId {
			maxId = exercise.Id
		}
	}

	// add +1 to len as IDs from the DB are 1 indexed, while slices are 0 indexed,
	// so we need to shift by one to have space for the latest row
	h.rows = make([]*ExerciseRow, maxId+1)

	for _, exercise := range exercises {
		h.rows[exercise.Id] = &ExerciseRow{exercise: exercise, parent: h}
	}
}
