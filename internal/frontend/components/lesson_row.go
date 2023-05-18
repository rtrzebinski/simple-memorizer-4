package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"net/url"
	"strconv"
)

type LessonRow struct {
	app.Compo
	parent *Lessons
	lesson models.Lesson
}

// The Render method is where the component appearance is defined.
func (c *LessonRow) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Td().Style("border", "1px solid black").Body(
			app.Text(c.lesson.Id),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(c.lesson.Name),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Button().Text("Edit").OnClick(c.onEdit(), fmt.Sprintf("%p", c)),
			app.Button().Text("Delete").OnClick(c.onDelete(c.lesson.Id), fmt.Sprintf("%p", c)),
			app.Button().Text("Exercises").OnClick(c.onExercises(c.lesson.Id), fmt.Sprintf("%p", c)),
		),
	)
}

func (c *LessonRow) onEdit() app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		c.parent.inputId = c.lesson.Id
		c.parent.inputName = c.lesson.Name
	}
}

func (c *LessonRow) onDelete(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		// delete lesson via API
		err := c.parent.api.DeleteLesson(models.Lesson{Id: id})
		if err != nil {
			app.Log(fmt.Errorf("failed to delete lesson: %w", err))
		}
		// create a new rows slice to be replaced in parent component
		rows := make([]*LessonRow, len(c.parent.rows))
		for i, row := range c.parent.rows {
			// add all rows but current one (which is being deleted)
			if i != id && row != nil {
				rows[i] = row
			}
		}
		// replace parent rows slice with a new one - this will update the UI
		c.parent.rows = rows
	}
}

func (c *LessonRow) onExercises(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		u, _ := url.Parse(pathExercises)

		// set lesson_id in the url
		params := u.Query()
		params.Add("lesson_id", strconv.Itoa(id))
		u.RawQuery = params.Encode()

		ctx.NavigateTo(u)
	}
}