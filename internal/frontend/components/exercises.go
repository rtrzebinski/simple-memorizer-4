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

	// add new exercise
	inputQuestion      string
	inputAnswer        string
	saveButtonDisabled bool
}

// The Render method is where the component appearance is defined.
func (h *Exercises) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.Div().Body(
			app.H2().Text("Add new exercise"),
			app.Textarea().ID("input_question").Cols(30).Rows(3).Placeholder("Question").
				Required(true).OnInput(h.ValueTo(&h.inputQuestion)).Text(h.inputQuestion),
			app.Br(),
			app.Textarea().ID("input_answer").Cols(30).Rows(3).Placeholder("Answer").
				Required(true).OnInput(h.ValueTo(&h.inputAnswer)).Text(h.inputAnswer),
			app.Br(),
			app.Button().Text("Save").OnClick(h.storeExercise).Disabled(h.saveButtonDisabled),
		),
		app.Div().Body(
			app.H2().Text("All exercises"),
			app.Table().Style("border", "1px solid black").Body(
				&ExerciseHeader{},
				app.Range(h.rows).Slice(func(i int) app.UI {
					if h.rows[i] != nil {
						app.Log("rendering", i, h.rows[i].exercise.Id, h.rows[i].exercise.Question)
						return h.rows[i].Render()
					}
					return nil
				}),
			),
		),
	)
}

func (h *Exercises) storeExercise(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// todo implement input validation
	if h.inputQuestion == "" || h.inputAnswer == "" {
		app.Log("invalid input")
		return
	}

	h.saveButtonDisabled = true

	err := h.api.StoreExercise(models.Exercise{
		Question: h.inputQuestion,
		Answer:   h.inputAnswer,
	})
	if err != nil {
		app.Log(fmt.Errorf("failed to store exercise: %w", err))
	}

	h.inputQuestion = ""
	h.inputAnswer = ""

	h.fetchAllExercises()

	h.saveButtonDisabled = false
}

// The OnMount method is run once component is mounted
func (h *Exercises) OnMount(ctx app.Context) {
	url := app.Window().URL()
	h.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
	h.fetchAllExercises()
}

func (h *Exercises) fetchAllExercises() {
	exercises, err := h.api.FetchAllExercises()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch all exercises: %w", err))
	}
	app.Log("fetchAllExercises")
	h.rows = make([]*ExerciseRow, len(exercises)+1)

	for _, exercise := range exercises {
		h.rows[exercise.Id] = &ExerciseRow{exercise: exercise, parent: h}
	}
}
