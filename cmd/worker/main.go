package main

import (
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
	"github.com/rtrzebinski/simple-memorizer-4/internal/worker"
	"github.com/rtrzebinski/simple-memorizer-4/internal/worker/cloudevents"
	"github.com/rtrzebinski/simple-memorizer-4/internal/worker/result"
	"github.com/rtrzebinski/simple-memorizer-4/internal/worker/storage/postgres"
	probes "github.com/rtrzebinski/simple-memorizer-4/pkg/probes"
	"github.com/rtrzebinski/simple-memorizer-4/pkg/signal"
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("service", "worker")

	if err := run(ctx, logger); err != nil {
		logger.Error(err.Error())
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	logger.Info("application starting")

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

	reader := postgres.NewReader(db)
	writer := postgres.NewWriter(db)
	service := worker.NewService(reader, writer, result.NewProjectionBuilder())
	handler := cloudevents.NewHandler(service)

	receiver := func(ctx context.Context, ev event.Event) error {
		return handler.Handle(ctx, ev)
	}

	// Start CloudEvents receiver and send errors to the channel
	go func() {
		logger.Info("initializing CloudEvents receiver")
		serverErrors <- ceClient.StartReceiver(ctx, receiver)
	}()

	// =========================================
	// Start probes
	// =========================================

	probeServer := probes.SetupProbeServer(cfg.ProbeAddr, db)

	// Start probe server and send errors to the channel
	go func() {
		logger.Info("initializing probe server", "addr", cfg.ProbeAddr)
		serverErrors <- probeServer.ListenAndServe()
	}()

	logger.Info("application running")

	// =========================================
	// Blocking main and waiting for shutdown.
	// =========================================

	done := signal.NewNotifier(ctx)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-done.Done():
		logger.Info("start shutdown")

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

	logger.Info("application completed")

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
