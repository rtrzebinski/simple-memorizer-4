package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const PathAuthRegister = "/register"

// Register is a component that displays the registration form
type Register struct {
	app.Compo
}

// NewRegister creates a new Register component
func NewRegister() *Register {
	return &Register{}
}

// The Render method is where the component appearance is defined.
func (compo *Register) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Text("Register page"),
		),
	)
}
