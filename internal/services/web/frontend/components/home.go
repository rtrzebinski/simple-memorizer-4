package components

import (
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/auth"
)

const PathHome = "/"

// Home is a component that displays the home page
type Home struct {
	app.Compo
	c APIClient

	// component vars
	showHome bool
	user     *frontend.UserProfile
}

// NewHome creates a new Home component
func NewHome(c APIClient) *Home {
	return &Home{
		user: &frontend.UserProfile{},
		c:    c,
	}
}

// The Render method is where the component appearance is defined.
func (compo *Home) Render() app.UI {
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
		app.Text("Welcome "+compo.user.Name),
		app.Br(),
		app.P().Body(
			app.Text("Home page"),
		),
	)
}

// The OnMount method is run once component is mounted
func (compo *Home) OnMount(ctx app.Context) {
	// auth check
	compo.user = auth.CheckUser(ctx)
	if compo.user == nil {
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
		return
	}
}
