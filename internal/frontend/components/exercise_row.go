package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"net/http"
)

type ExerciseRow struct {
	app.Compo
	api *frontend.ApiClient

	parent   *Exercises
	exercise models.Exercise
}

// The Render method is where the component appearance is defined.
func (h *ExerciseRow) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.Id),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.Question),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.Answer),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.BadAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.GoodAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Button().Text("Delete").OnClick(func(ctx app.Context, e app.Event) {
				app.Logf("delete %d", h.exercise.Id)
				app.Logf("delete %d", h.exercise.Question)
				err := h.api.DeleteExercise(h.exercise)
				if err != nil {
					app.Log(fmt.Errorf("failed to delete exercise: %w", err))
				}
				// remove deleted row from the parent slice
				rows := make([]*ExerciseRow, len(h.parent.rows))
				for i, row := range h.parent.rows {
					if i != h.exercise.Id && row != nil {
						app.Log("inserting", i, row.exercise.Id, row.exercise.Question)
						rows[i] = row
					}
				}
				app.Log(h.parent.rows)
				h.parent.rows = rows
				app.Log(rows)
			}),
		),
	)
}

// The OnMount method is run once component is mounted
func (h *ExerciseRow) OnMount(ctx app.Context) {
	url := app.Window().URL()
	h.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
}
