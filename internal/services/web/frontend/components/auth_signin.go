package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const PathAuthSignIn = "/sign-in"

// SignIn is a component that displays the sign-in form
type SignIn struct {
	app.Compo
	c APIClient
}

// NewSignIn creates a new sign-in component
func NewSignIn(c APIClient) *SignIn {
	return &SignIn{c: c}
}

// The Render method is where the component appearance is defined.
func (compo *SignIn) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Text("SignIn page"),
		),
	)
}
