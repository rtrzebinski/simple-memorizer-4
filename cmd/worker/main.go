package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	cepubsub "github.com/cloudevents/sdk-go/protocol/pubsub/v2"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	probes "github.com/rtrzebinski/simple-memorizer-4/internal/probes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker/pubsub"
	"github.com/rtrzebinski/simple-memorizer-4/internal/signal"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
)

type config struct {
	Db struct {
		Driver string `envconfig:"DB_DRIVER" default:"postgres"`
		DSN    string `envconfig:"DB_DSN" default:"postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable&timezone=Europe/Warsaw"`
	}
	PubSub struct {
		ProjectID       string   `envconfig:"PUBSUB_PROJECT_ID" default:"project-dev"`
		TopicID         string   `envconfig:"PUBSUB_TOPIC_ID" default:"topic-dev"`
		SubscriptionIDs []string `envconfig:"PUBSUB_SUBSCRIPTION_IDS" default:"subscription-dev"`
	}
	ProbeAddr       string        `envconfig:"PROBE_ADDRESS" default:"0.0.0.0:9091"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		slog.Error(err.Error())
	}
}

func run(ctx context.Context) error {
	slog.Info("application starting", "service", "worker")

	// Version
	var version string
	file, err := os.Open("version")
	if err == nil {
		defer func() {
			err := file.Close()
			if err != nil {
				slog.Warn("failed to close file", "error", err, "service", "worker")
			}
		}()

		scanner := bufio.NewScanner(file)

		if scanner.Scan() {
			version = scanner.Text()
			slog.Info("version", "version", version, "service", "worker")
		} else {
			slog.Info("version unknown", "service", "worker")
		}
	}

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
		cfg.PubSub.SubscriptionIDs,
	)
	if err != nil {
		return err
	}

	// Make a channel to listen for errors.
	// Use a buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// =========================================
	// Start worker
	// =========================================

	reader := postgres.NewWorkerReader(db)
	writer := postgres.NewWorkerWriter(db)
	service := worker.NewService(reader, writer)
	handler := pubsub.NewHandler(service)

	receiver := func(ctx context.Context, ev event.Event) error {
		return handler.Handle(ctx, ev)
	}

	// Start CloudEvents receiver and send errors to the channel
	go func() {
		slog.Info("initializing CloudEvents receiver", "service", "worker")
		serverErrors <- ceClient.StartReceiver(ctx, receiver)
	}()

	// =========================================
	// Start probes
	// =========================================

	probeServer := probes.SetupProbeServer(cfg.ProbeAddr, db, nil)

	// Start probe server and send errors to the channel
	go func() {
		slog.Info("initializing probe server", "addr", cfg.ProbeAddr, "service", "worker")
		serverErrors <- probeServer.ListenAndServe()
	}()

	slog.Info("application running", "service", "worker")

	// =========================================
	// Blocking main and waiting for shutdown.
	// =========================================

	done := signal.NewNotifier(ctx)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-done.Done():
		slog.Info("start shutdown", "service", "worker")

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

	slog.Info("application completed", "service", "worker")

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
