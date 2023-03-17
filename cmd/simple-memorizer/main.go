package main

import (
	"context"
	"database/sql"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/client/components"
	"github.com/rtrzebinski/simple-memorizer-4/internal/server"
	"github.com/rtrzebinski/simple-memorizer-4/internal/server/storage/postgres"
	"log"
	"net/http"
)

type config struct {
	Db struct {
		Driver string `envconfig:"DB_DRIVER" default:"postgres"`
		DSN    string `envconfig:"DB_DSN" default:"postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable"`
	}
	Api struct {
		Port string `envconfig:"API_PORT" default:":8000"`
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		panic(err)
	}
}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func run(ctx context.Context) error {
	log.Println("App started..")

	// The first thing to do is to associate the home component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", &components.Home{})

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for pre compiling
	// instructions.
	app.RunWhenOnBrowser()

	// Configuration
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return err
	}

	// Database connection
	db, err := sql.Open(cfg.Db.Driver, cfg.Db.DSN)
	if err != nil {
		return err
	}

	// Dependencies
	r := postgres.NewReader(db)
	w := postgres.NewWriter(db)

	// Handle home page
	http.Handle("/", &app.Handler{
		Name:        "Home",
		Description: "Home page",
	})

	// Start API server
	if err = server.ListenAndServe(r, w, cfg.Api.Port); err != nil {
		return err
	}

	return nil
}
