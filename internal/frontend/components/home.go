package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"net/http"
)

var pathHome = "/"

// A Home component
type Home struct {
	app.Compo
	api *frontend.ApiClient
}

// The OnMount method is run once component is mounted
func (c *Home) OnMount(ctx app.Context) {
	url := app.Window().URL()
	c.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
}

// The Render method is where the component appearance is defined.
func (c *Home) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Text("Home page"),
		),
	)
}
