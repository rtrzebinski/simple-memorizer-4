package components

import (
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"strconv"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/csv"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/validation"
)

const PathExercises = "/exercises"

// Exercises is a component that displays all exercises of a lesson
type Exercises struct {
	app.Compo
	c APIClient

	// component vars
	lesson frontend.Lesson
	rows   []*ExerciseRow
	user   *frontend.User

	// upsert exercise form
	formVisible        bool
	validationErrors   []error
	inputId            int
	inputQuestion      string
	inputAnswer        string
	saveButtonDisabled bool
}

// NewExercises creates a new Exercises component
func NewExercises(c APIClient) *Exercises {
	return &Exercises{
		c:    c,
		user: &frontend.User{},
	}
}

// The OnMount method is run once component is mounted
func (compo *Exercises) OnMount(ctx app.Context) {
	// auth check
	user, err := auth.User(ctx)
	if err != nil {
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}
	compo.user = user

	u := app.Window().URL()

	// find lesson_id in the url
	lessonId, err := strconv.Atoi(u.Query().Get("lesson_id"))
	if err != nil {
		app.Log(fmt.Errorf("failed to convert lesson_id: %w", err))
		return
	}

	compo.lesson = frontend.Lesson{Id: lessonId}

	compo.hydrateLesson(ctx)
	compo.displayExercisesOfLesson(ctx)
}

// hydrateLesson fetch lesson details
func (compo *Exercises) hydrateLesson(ctx app.Context) {
	accessToken, err := auth.Token(ctx)
	if err != nil {
		slog.Error("failed to get token", "err", err)
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}
	err = compo.c.HydrateLesson(ctx, &compo.lesson, accessToken)
	if err != nil {
		app.Log(fmt.Errorf("failed to hydrate lesson: %w", err))
	}
}

// The Render method is where the component appearance is defined.
func (compo *Exercises) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.P().Body(
				app.A().Href(PathHome).Text("Home"),
				app.Text(" | "),
				app.A().Href(PathLessons).Text("Lessons"),
				app.Text(" | "),
				app.A().Href(PathAuthLogout).Text("Logout"),
				app.Text(" | "),
				app.Text(app.Getenv("version")),
			),
		),
		app.P().Body(
			app.Button().Text("Start learning").OnClick(compo.handleStartLearning).Disabled(compo.lesson.ExerciseCount < 2),
			app.Button().Text("Add a new exercise").OnClick(compo.handleAddExercise).Hidden(compo.formVisible),
			app.A().Href(http.ExportLessonCsv+"?lesson_id="+strconv.Itoa(compo.lesson.Id)).Download("").Body(
				app.Button().Text("CSV export"),
			),
			//app.Label().For("csv-upload-button").Text("CSV import"), todo style with bootstrap
			app.Input().ID("csv-upload-button").Type("file").Name("csv-import.csv").Accept(".csv").OnInput(compo.handleCsvUpload),
		),
		app.P().Body(
			app.Text("Lesson name: "),
			app.Text(compo.lesson.Name),
		),
		app.P().Body(
			app.Text("Exercises: "),
			app.Text(compo.lesson.ExerciseCount),
		),
		app.Div().Body(
			app.H3().Text("Add a new exercise"),
			app.P().Body(
				app.Range(compo.validationErrors).Slice(func(i int) app.UI {
					return app.Div().Body(
						app.Text(compo.validationErrors[i].Error()),
						app.Br(),
					).Style("color", "red")
				}),
			),
			app.Textarea().ID("input_question").Cols(30).Rows(3).Placeholder("Question").
				Required(true).OnInput(compo.ValueTo(&compo.inputQuestion)).Text(compo.inputQuestion),
			app.Br(),
			app.Textarea().ID("input_answer").Cols(30).Rows(3).Placeholder("Answer").
				Required(true).OnInput(compo.ValueTo(&compo.inputAnswer)).Text(compo.inputAnswer),
			app.Br(),
			app.Button().Text("Save").OnClick(compo.handleSave).Disabled(compo.saveButtonDisabled),
			app.Button().Text("Cancel").OnClick(compo.handleCancel),
		).Hidden(!compo.formVisible),
		app.Div().Body(
			app.H3().Text("Exercises"),
			app.Table().Style("border", "1px solid black").Body(
				&ExerciseHeader{},
				app.Range(compo.rows).Slice(func(i int) app.UI {
					if compo.rows[i] != nil {
						return compo.rows[i].Render()
					}
					return nil
				}),
			),
		),
	)
}

// handleStartLearning start learning a current lesson
func (compo *Exercises) handleStartLearning(ctx app.Context, _ app.Event) {
	u, _ := url.Parse(PathLearn)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(compo.lesson.Id))
	u.RawQuery = params.Encode()

	ctx.NavigateTo(u)
}

// handleAddExercise display add exercise form
func (compo *Exercises) handleAddExercise(_ app.Context, _ app.Event) {
	compo.formVisible = true
}

// handleCsvUpload upload exercises from a CSV file
func (compo *Exercises) handleCsvUpload(ctx app.Context, e app.Event) {
	file := e.Get("target").Get("files").Index(0)

	// validate file type which can be changed by the user in the file picker
	if file.Get("type").String() != "text/csv" {
		app.Log("invalid file type", file.Get("type").String())
		return
	}

	// read bytes from uploaded file
	data, err := readFile(file)
	if err != nil {
		app.Log(err)
		return
	}

	// extract records from bytes slice
	records, err := csv.ReadAll(data)
	if err != nil {
		app.Log(err)
		return
	}

	// prepare exercises to be stored
	var exercises []frontend.Exercise
	for _, rec := range records {
		exercises = append(exercises, frontend.Exercise{
			Lesson:   &compo.lesson,
			Question: rec[0],
			Answer:   rec[1],
		})
	}

	// store uploaded exercises
	accessToken, err := auth.Token(ctx)
	if err != nil {
		slog.Error("failed to get token", "err", err)
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}
	err = compo.c.StoreExercises(ctx, exercises, accessToken)
	if err != nil {
		app.Log(err)
		return
	}

	// reset file input for next upload
	e.Get("target").Set("value", "")

	// reload the UI
	compo.hydrateLesson(ctx)
	compo.displayExercisesOfLesson(ctx)
}

// readFile some JS magic converting uploaded file to a slice of bytes
// https://github.com/maxence-charriere/go-app/issues/882
func readFile(file app.Value) (data []byte, err error) {
	done := make(chan bool)

	// https://developer.mozilla.org/en-US/docs/Web/API/FileReader
	reader := app.Window().Get("FileReader").New()
	reader.Set("onloadend", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		done <- true
		return nil
	}))
	reader.Call("readAsArrayBuffer", file)
	<-done

	readerError := reader.Get("error")

	if !readerError.IsNull() {
		err = fmt.Errorf("file reader error : %p", readerError.Get("message"))
	} else {
		uint8Array := app.Window().Get("Uint8Array").New(reader.Get("result"))
		data = make([]byte, uint8Array.Length())
		app.CopyBytesToGo(data, uint8Array)
	}

	return data, err
}

// handleSave create new or update existing exercise
func (compo *Exercises) handleSave(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// exercise to be saved
	exercise := frontend.Exercise{
		Id:       compo.inputId,
		Question: compo.inputQuestion,
		Answer:   compo.inputAnswer,
		Lesson: &frontend.Lesson{
			Id: compo.lesson.Id,
		},
	}

	// extract other questions to validate against
	var questions []string
	for i, row := range compo.rows {
		if row != nil && compo.rows[i].exercise.Id != compo.inputId {
			questions = append(questions, compo.rows[i].exercise.Question)
		}
	}

	// validate input
	validator := validation.ValidateUpsertExercise(exercise, questions)
	if validator.Failed() {
		compo.validationErrors = validator.Errors()

		return
	}

	// disable submit button to avoid duplicated requests
	compo.saveButtonDisabled = true

	// save exercise
	accessToken, err := auth.Token(ctx)
	if err != nil {
		slog.Error("failed to get token", "err", err)
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}

	err = compo.c.UpsertExercise(ctx, exercise, accessToken)
	if err != nil {
		app.Log(fmt.Errorf("failed to save exercise: %w", err))
	}

	// hide the input form on row edit, but keep open on adding new
	// because it is common to add a few exercises one after another
	if compo.inputId != 0 {
		compo.formVisible = false
	}

	// reset form
	compo.resetForm()

	// reload the UI
	compo.hydrateLesson(ctx)
	compo.displayExercisesOfLesson(ctx)
}

// handleCancel handle cancel button click
func (compo *Exercises) handleCancel(_ app.Context, _ app.Event) {
	compo.resetForm()
	compo.formVisible = false
}

// resetForm reset form fields
func (compo *Exercises) resetForm() {
	compo.inputId = 0
	compo.inputQuestion = ""
	compo.inputAnswer = ""
	compo.validationErrors = nil
	compo.saveButtonDisabled = false
}

// displayExercisesOfLesson fetch exercises and display them
func (compo *Exercises) displayExercisesOfLesson(ctx app.Context) {
	accessToken, err := auth.Token(ctx)
	if err != nil {
		slog.Error("failed to get token", "err", err)
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}
	oldestExerciseID := 1 // Set the oldest exercise ID to 1, as we are displaying all exercises
	exercises, err := compo.c.FetchExercises(ctx, compo.lesson, oldestExerciseID, accessToken)
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch exercises of lesson: %w", err))
	}

	// no entries
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
	compo.rows = make([]*ExerciseRow, maxId+1)

	for _, row := range exercises {
		compo.rows[row.Id] = &ExerciseRow{exercise: row, parent: compo}
	}

	// reverse the order to display the latest exercises first
	slices.Reverse(compo.rows)
}
