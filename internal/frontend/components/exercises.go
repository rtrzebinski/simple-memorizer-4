package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"net/http"
)

type Exercises struct {
	app.Compo
	api  *frontend.ApiClient
	rows []*ExerciseRow
}

// The Render method is where the component appearance is defined.
func (h *Exercises) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.Table().Style("border", "1px solid black").Body(
			&ExerciseHeader{},
			app.Range(h.rows).Slice(func(i int) app.UI {
				return h.rows[i]
			}),
		),
	)
}

// The OnMount method is run once component is mounted
func (h *Exercises) OnMount(ctx app.Context) {
	url := app.Window().URL()
	h.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)

	exercises, err := h.api.FetchExercises()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch exercises: %w", err))
	}

	for _, exercise := range exercises {
		h.rows = append(h.rows, &ExerciseRow{question: exercise.Question, answer: exercise.Answer})
	}
}
