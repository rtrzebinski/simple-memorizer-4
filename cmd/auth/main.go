package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	probes "github.com/rtrzebinski/simple-memorizer-4/pkg/probes"
	"github.com/rtrzebinski/simple-memorizer-4/pkg/signal"
)

type config struct {
	ProbeAddr       string        `envconfig:"PROBE_ADDRESS" default:"0.0.0.0:9092"`
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
	slog.Info("application starting", "service", "auth")

	// Version
	file, err := os.Open("version")
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var version string

	if scanner.Scan() {
		version = scanner.Text()
		slog.Info("version", "version", version, "service", "auth")

	} else {
		slog.Info("version unknown", "service", "auth")
	}

	// Configuration
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return err
	}

	// Make a channel to listen for errors.
	// Use a buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// =========================================
	// Start auth
	// =========================================

	// todo

	// =========================================
	// Start probes
	// =========================================

	// todo add a db check later
	probeServer := probes.SetupProbeServer(cfg.ProbeAddr, nil)

	// Start probe server and send errors to the channel
	go func() {
		slog.Info("initializing probe server", "addr", cfg.ProbeAddr, "service", "auth")
		serverErrors <- probeServer.ListenAndServe()
	}()

	slog.Info("application running", "service", "auth")

	// =========================================
	// Blocking main and waiting for shutdown.
	// =========================================

	done := signal.NewNotifier(ctx)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-done.Done():
		slog.Info("start shutdown", "service", "auth")

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

	slog.Info("application completed", "service", "auth")

	return nil
}
