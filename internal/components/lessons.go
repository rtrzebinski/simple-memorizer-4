package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/http/rest"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"net/http"
)

var pathLessons = "/lessons"

type Lessons struct {
	app.Compo
	reader *rest.Reader
	writer *rest.Writer
	rows   []*LessonRow

	addLessonFormVisible bool
	inputId              int
	inputName            string
	storeButtonDisabled  bool
}

// The OnMount method is run once component is mounted
func (c *Lessons) OnMount(ctx app.Context) {
	url := app.Window().URL()
	c.reader = rest.NewReader(&http.Client{}, url.Host, url.Scheme)
	c.writer = rest.NewWriter(&http.Client{}, url.Host, url.Scheme)
	c.displayAllLessons()
}

// The Render method is where the component appearance is defined.
func (c *Lessons) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Button().Text("Add a new lesson").OnClick(c.handleAddLesson).Hidden(c.addLessonFormVisible),
		),
		app.Div().Body(
			app.H3().Text("Add a new lesson"),
			app.Textarea().ID("input_name").Cols(30).Rows(3).Placeholder("Name").
				Required(true).OnInput(c.ValueTo(&c.inputName)).Text(c.inputName),
			app.Br(),
			app.Button().Text("Store").OnClick(c.handleStore).Disabled(c.storeButtonDisabled),
			app.Button().Text("Cancel").OnClick(c.handleCancel),
		).Hidden(!c.addLessonFormVisible),
		app.Div().Body(
			app.H3().Text("Lessons"),
			app.Table().Style("border", "1px solid black").Body(
				&LessonHeader{},
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

// handleAddLesson display add lesson form
func (c *Lessons) handleAddLesson(ctx app.Context, e app.Event) {
	c.addLessonFormVisible = true
}

// handleStore create new or update existing lesson
func (c *Lessons) handleStore(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// todo implement input validation
	if c.inputName == "" {
		app.Log("invalid input")
		return
	}

	c.storeButtonDisabled = true

	err := c.writer.StoreLesson(models.Lesson{
		Id:   c.inputId,
		Name: c.inputName,
	})
	if err != nil {
		app.Log(fmt.Errorf("failed to store lesson: %w", err))
	}

	c.inputId = 0
	c.inputName = ""

	c.displayAllLessons()

	c.storeButtonDisabled = false

	// hide the inout form on success
	c.addLessonFormVisible = false
}

// handleCancel cleanup input UI
func (c *Lessons) handleCancel(ctx app.Context, e app.Event) {
	c.inputId = 0
	c.inputName = ""
	c.addLessonFormVisible = false
}

func (c *Lessons) displayAllLessons() {
	lessons, err := c.reader.FetchAllLessons()
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch all lessons: %w", err))
	}

	// no entries in the database
	if len(lessons) == 0 {
		return
	}

	// find maxId so we know the rows slice capacity
	maxId := lessons[0].Id
	for _, row := range lessons {
		if row.Id > maxId {
			maxId = row.Id
		}
	}

	// add +1 to len as IDs from the DB are 1 indexed, while slices are 0 indexed,
	// so we need to shift by one to have space for the latest row
	c.rows = make([]*LessonRow, maxId+1)

	for _, row := range lessons {
		c.rows[row.Id] = &LessonRow{lesson: row, parent: c}
	}
}