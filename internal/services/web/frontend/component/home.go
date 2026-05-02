package component

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/component/auth"
)

const PathHome = "/"

// Home is a component that displays the home page
type Home struct {
	app.Compo
	c APIClient

	// component vars
	showHome    bool
	userProfile *frontend.UserProfile
}

// NewHome creates a new Home component
func NewHome(c APIClient) *Home {
	return &Home{
		userProfile: &frontend.UserProfile{},
		c:           c,
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
		app.Text("Welcome "+compo.userProfile.Name),
		app.Br(),
		app.P().Body(
			app.Text("Home page"),
		),
	)
}

// The OnMount method is run once component is mounted
func (compo *Home) OnMount(ctx app.Context) {
	// auth check
	compo.userProfile = auth.GetUser(ctx)
	if compo.userProfile == nil {
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
		return
	}

	compo.displayProfile(ctx)
}

func (compo *Home) displayProfile(ctx app.Context) {
	userProfile, err := compo.c.UserProfile(ctx)
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch userProfile profile: %w", err))
		if errors.Is(err, frontend.ErrUnauthorized) {
			auth.DelUser(ctx)
			ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
		}
		return
	}

	compo.userProfile = userProfile
}
