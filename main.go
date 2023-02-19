package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"log"
	"math/rand"
	"net/http"
)

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// The first thing to do is to associate the hello component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", &hello{})

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for precompiling
	// instructions.
	app.RunWhenOnBrowser()

	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

// hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type hello struct {
	app.Compo

	question    string
	answer      string
	showAnswer  bool
	goodAnswers int
	badAnswers  int
}

// init fetch random exercises, showAnswer question, hide answer
func (h *hello) init() {
	h.showAnswer = false
	question, answer := exercises()
	h.question = question
	h.answer = answer
}

// The Render method is where the component appearance is defined.
func (h *hello) Render() app.UI {
	// will be re-run on every button click
	if h.question == "" && h.answer == "" {
		h.init()
	}

	return app.Div().Body(
		app.P().Body(
			app.Button().
				Text("Next exercise").
				OnClick(func(ctx app.Context, e app.Event) {
					h.init()
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("Show Answer").
				OnClick(func(ctx app.Context, e app.Event) {
					h.showAnswer = true
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("Good Answer").
				OnClick(func(ctx app.Context, e app.Event) {
					h.goodAnswers++
					h.init()
				}).
				Style("margin-right", "10px"),
			app.Button().
				Text("Bad Answer").
				OnClick(func(ctx app.Context, e app.Event) {
					h.badAnswers++
					h.init()
				}).
				Style("margin-right", "10px"),
		),
		app.H2().Body(
			app.Text("What is the capital of "),
			app.If(h.question != "",
				app.Text(h.question),
			).Else(
				app.Text("..."),
			),
			app.Text("?"),
		),
		app.H2().Body(
			app.If(h.answer != "",
				app.Text(h.answer),
			).Else(
				app.Text("?"),
			),
		).Hidden(!h.showAnswer),
		app.P().Body(
			app.Text("Good answers: "),
			app.Text(h.goodAnswers),
		),
		app.P().Body(
			app.Text("Bad answers: "),
			app.Text(h.badAnswers),
		),
	)
}

func exercises() (string, string) {
	capitals := map[string]string{
		"Poland":      "Warsaw",
		"Germany":     "Berlin",
		"France":      "Paris",
		"Netherlands": "Amsterdam",
		"Spain":       "Madrid",
	}

	i := rand.Intn(len(capitals))
	for k := range capitals {
		if i == 0 {
			return k, capitals[k]
		}
		i--
	}

	panic("never")
}
