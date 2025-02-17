package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	bhttp "github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	cepubsub "github.com/cloudevents/sdk-go/protocol/pubsub/v2"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	probes "github.com/rtrzebinski/simple-memorizer-4/internal/probes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/cloudevents"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components"
	fhttp "github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/signal"
	storage "github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres/web"
)

type config struct {
	Web struct {
		Port     string `envconfig:"WEB_PORT" default:":8000"`
		CertFile string `envconfig:"WEB_CERT_FILE" default:"ssl/localhost-cert.pem"`
		KeyFile  string `envconfig:"WEB_KEY_FILE" default:"ssl/localhost-key.pem"`
	}
	Db struct {
		Driver string `envconfig:"DB_DRIVER" default:"postgres"`
		DSN    string `envconfig:"DB_DSN" default:"postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable&timezone=Europe/Warsaw"`
	}
	PubSub struct {
		ProjectID string `envconfig:"PUBSUB_PROJECT_ID" default:"project-dev"`
		TopicID   string `envconfig:"PUBSUB_TOPIC_ID" default:"topic-dev"`
	}
	ProbeAddr       string        `envconfig:"PROBE_ADDRESS" default:"0.0.0.0:9090"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		slog.Error(err.Error())
	}
}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func run(ctx context.Context) error {
	slog.Info("application starting", "service", "web")

	u := app.Window().URL()

	// create a service to be injected into components
	apiClient := fhttp.NewClient(fhttp.NewHTTPCaller(&http.Client{}, u.Host, u.Scheme))

	// The first thing to do is to associate the home component with a path.
	//
	// This is done by calling the Route() function, which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route(components.PathHome, func() app.Composer { return components.NewHome() })

	// Associate other frontend routes
	app.Route(components.PathLessons, func() app.Composer { return components.NewLessons(apiClient) })
	app.Route(components.PathExercises, func() app.Composer { return components.NewExercises(apiClient) })
	app.Route(components.PathLearn, func() app.Composer { return components.NewLearn(apiClient) })

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

	// Version
	var version string
	file, err := os.Open("version")
	if err == nil {
		defer func() {
			err := file.Close()
			if err != nil {
				slog.Warn("failed to close file", "error", err, "service", "web")
			}
		}()

		scanner := bufio.NewScanner(file)

		if scanner.Scan() {
			version = scanner.Text()
			slog.Info("version", "version", version, "service", "web")
		} else {
			slog.Info("version unknown", "service", "web")
		}
	}

	// Handle home page
	http.Handle("/", &app.Handler{
		Name:        "Home",
		Description: "Home page",
		Icon: app.Icon{
			Default:  "/web/logo-192.png",
			Large:    "/web/logo-512.png",
			Maskable: "/web/logo-192.png",
		},
		Scripts: []string{
			"/web/swiped-events.js",
		},
		Styles: []string{
			// todo find a way to only load on a learning page
			//"/web/hello.css",
		},
		Env: map[string]string{
			"version": version,
		},
	})

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

	// CloudEvents client
	ceClient, err := createCloudEventsClient(ctx,
		cfg.PubSub.ProjectID,
		cfg.PubSub.TopicID,
		[]string{},
	)
	if err != nil {
		return err
	}

	// Make a channel to listen for errors.
	// Use a buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// =========================================
	// Start server
	// =========================================

	reader := storage.NewReader(db)
	writer := storage.NewWriter(db)
	publisher := cloudevents.NewPublisher(ceClient)
	service := backend.NewService(reader, writer, publisher)

	go func() {
		slog.Info("initializing server", "port", cfg.Web.Port, "service", "web")
		serverErrors <- bhttp.ListenAndServe(service, cfg.Web.Port, cfg.Web.CertFile, cfg.Web.KeyFile)
	}()

	// =========================================
	// Start probes
	// =========================================

	probeServer := probes.SetupProbeServer(cfg.ProbeAddr, db)

	// Start probe server and send errors to the channel
	go func() {
		slog.Info("initializing probe server", "addr", cfg.ProbeAddr, "service", "web")
		serverErrors <- probeServer.ListenAndServe()
	}()

	slog.Info("application running", "service", "web")

	// =========================================
	// Blocking main and waiting for shutdown.
	// =========================================

	done := signal.NewNotifier(ctx)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-done.Done():
		slog.Info("start shutdown", "service", "web")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
		defer cancel()

		// Shutdown gracefully on signal received
		if err := probeServer.Shutdown(ctx); err != nil {
			log.Print(fmt.Errorf("failed to gracefully shutdown the probe server %w", err))

			if err = probeServer.Close(); err != nil {
				return fmt.Errorf("could not stop probe server gracefully: %w", err)
			}
		}
	}

	slog.Info("application completed", "service", "web")

	return nil
}

func createCloudEventsClient(
	ctx context.Context,
	projectID string,
	topicID string,
	subscriptionIDs []string,
) (client.Client, error) {
	opts := make([]cepubsub.Option, 0)

	for _, sID := range subscriptionIDs {
		opts = append(opts, cepubsub.WithSubscriptionID(strings.Trim(sID, "\n")))
	}

	opts = append(opts,
		cepubsub.WithProjectID(projectID),
		cepubsub.WithTopicID(topicID),
		cepubsub.AllowCreateTopic(true),
	)

	t, err := cepubsub.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create transport for topic id (%q): %w", topicID, err)
	}

	c, err := ce.NewClient(t, ce.WithTimeNow(), ce.WithUUIDs())
	if err != nil {
		return nil, fmt.Errorf("create client for topic id (%q): %w", topicID, err)
	}

	return c, nil
}
