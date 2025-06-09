package components

import (
	"fmt"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/auth"
)

// ExerciseRow is a component that displays a single row of the exercises table
type ExerciseRow struct {
	app.Compo
	parent   *Exercises
	exercise frontend.Exercise
}

// The Render method is where the component appearance is defined.
func (compo *ExerciseRow) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.Question),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.Answer),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.BadAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.GoodAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.GoodAnswersPercent()),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Button().Text("Delete").OnClick(compo.onDelete(compo.exercise.Id), app.EventScope(fmt.Sprintf("%p", compo))),
			app.Button().Text("Edit").OnClick(compo.onEdit(), app.EventScope(fmt.Sprintf("%p", compo))),
			app.Button().Text("LearnSince").OnClick(compo.onLearnSince(), app.EventScope(fmt.Sprintf("%p", compo))),
		),
	)
}

// onDelete handles delete button click
func (compo *ExerciseRow) onDelete(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		// confirmation dialog
		if !app.Window().Call("confirm", "Delete a record?").Bool() {
			return
		}

		accessToken, err := auth.Token(ctx)
		if err != nil {
			slog.Error("failed to get token", "err", err)
			ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
		}
		// delete exercise
		err = compo.parent.c.DeleteExercise(ctx, frontend.Exercise{Id: id}, accessToken)
		if err != nil {
			app.Log(fmt.Errorf("failed to delete exercise: %w", err))
		}
		// create a new rows slice to be replaced in parent component
		rows := make([]*ExerciseRow, 0, len(compo.parent.rows))
		for _, row := range compo.parent.rows {
			if row != nil && row.exercise.Id != id {
				rows = append(rows, row)
			}
		}
		// replace parent rows slice with a new one - this will update the UI
		compo.parent.rows = rows
		compo.parent.hydrateLesson(ctx)
	}
}

// onEdit handles edit button click
func (compo *ExerciseRow) onEdit() app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		compo.parent.inputId = compo.exercise.Id
		compo.parent.inputQuestion = compo.exercise.Question
		compo.parent.inputAnswer = compo.exercise.Answer
		compo.parent.formVisible = true
		compo.parent.validationErrors = nil
	}
}

// onLearnSince handles LearnSince button click
// It navigates to the learn page with the oldest exercise ID and lesson ID as parameters.
// This allows the user to start learning from the latest exercise in the lesson, for instance only latest X exercises.
func (compo *ExerciseRow) onLearnSince() app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		u, _ := url.Parse(PathLearn)
		params := u.Query()
		params.Add("lesson_id", strconv.Itoa(compo.parent.lesson.Id))
		params.Add("oldest_exercise_id", strconv.Itoa(compo.exercise.Id))
		u.RawQuery = params.Encode()
		ctx.NavigateTo(u)
	}
}
