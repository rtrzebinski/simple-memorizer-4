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

// LessonRow is a component that displays a row in the lessons table
type LessonRow struct {
	app.Compo
	parent *Lessons
	lesson frontend.Lesson
}

// The Render method is where the component appearance is defined.
func (compo *LessonRow) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.lesson.Id),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.lesson.Name),
		),
		//app.Td().Style("border", "1px solid black").Body(
		//	app.Text(compo.lesson.Description),
		//),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.lesson.ExerciseCount),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Button().Text("Edit").OnClick(compo.onEdit(), app.EventScope(fmt.Sprintf("%p", compo))),
			app.Button().Text("Delete").OnClick(compo.onDelete(compo.lesson.Id), app.EventScope(fmt.Sprintf("%p", compo))),
			app.Button().Text("Exercises").OnClick(compo.onExercises(compo.lesson.Id), app.EventScope(fmt.Sprintf("%p", compo))),
			app.Button().Text("Learn").OnClick(compo.onLearn(compo.lesson.Id), app.EventScope(fmt.Sprintf("%p", compo))).
				// Learning empty lessons not allowed
				Disabled(compo.lesson.ExerciseCount < 2),
		),
	)
}

// onEdit handles edit button click
func (compo *LessonRow) onEdit() app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		compo.parent.inputId = compo.lesson.Id
		compo.parent.inputName = compo.lesson.Name
		compo.parent.inputDescription = compo.lesson.Description
		compo.parent.formVisible = true
		compo.parent.validationErrors = nil
	}
}

// onDelete handles delete button click
func (compo *LessonRow) onDelete(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		// delete lesson
		accessToken, err := auth.Token(ctx)
		if err != nil {
			slog.Error("failed to get token", "err", err)
			ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
		}
		err = compo.parent.c.DeleteLesson(ctx, frontend.Lesson{Id: id}, accessToken)
		if err != nil {
			app.Log(fmt.Errorf("failed to delete lesson: %w", err))
		}
		// create a new rows slice to be replaced in parent component
		rows := make([]*LessonRow, len(compo.parent.rows))
		for i, row := range compo.parent.rows {
			// add all rows but current one (which is being deleted)
			if i != id && row != nil {
				rows[i] = row
			}
		}
		// replace parent rows slice with a new one - this will update the UI
		compo.parent.rows = rows
	}
}

// onExercises handles exercises button click
func (compo *LessonRow) onExercises(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		u, _ := url.Parse(PathExercises)

		// set lesson_id in the url
		params := u.Query()
		params.Add("lesson_id", strconv.Itoa(id))
		u.RawQuery = params.Encode()

		ctx.NavigateTo(u)
	}
}

// onLearn handles learn button click
func (compo *LessonRow) onLearn(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		u, _ := url.Parse(PathLearn)

		// set lesson_id in the url
		params := u.Query()
		params.Add("lesson_id", strconv.Itoa(id))
		u.RawQuery = params.Encode()

		ctx.NavigateTo(u)
	}
}
