package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Navigation is a component that displays the navigation bar
type Navigation struct {
	app.Compo

	showHome    bool
	showLessons bool
}

// NewNavigation creates a new Navigation component
func NewNavigation() *Navigation {
	return &Navigation{}
}

// The Render method is where the component appearance is defined.
func (compo *Navigation) Render() app.UI {
	return app.Div().Body(
		app.P().Body(
			// todo for some reason this does not work:
			//app.If(compo.showHome == true, func() app.UI {
			//	return app.A().Href(PathHome).Text("Home")
			//}),
			//app.A().Href(PathHome).Text("Home | ").Hidden(!compo.showHome),
			//app.A().Href(PathLessons).Text("Lessons | ").Hidden(!compo.showLessons),
			//app.A().Href(PathAuthRegister).Text("Register | "),
			//app.A().Href(PathAuthSignIn).Text("SignIn | "),
			//app.A().Href(PathAuthLogout).Text("Logout | "),
			app.A().Href(PathHome).Text("Home"),
			app.Text(" | "),
			app.A().Href(PathLessons).Text("Lessons"),
			app.Text(" | "),
			app.A().Href(PathAuthRegister).Text("Register"),
			app.Text(" | "),
			app.A().Href(PathAuthSignIn).Text("SignIn"),
			app.Text(" | "),
			app.A().Href(PathAuthLogout).Text("Logout"),
			app.Text(" | "),
			app.Text(app.Getenv("version")),
		),
	)
}

/*
	app.If(compo.showHome == true, func() app.UI {
		return app.A().Href(PathHome).Text("Home")
	}),
	app.If(compo.showHome == true, func() app.UI {
		return app.Text(" | ")
	}),
	app.If(compo.showLessons == true, func() app.UI {
		return app.A().Href(PathLessons).Text("Lessons")
	}),
	app.If(compo.showLessons == true, func() app.UI {
		return app.Text(" | ")
	}),
*/
