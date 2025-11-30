package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/kelseyhightower/envconfig"
	gengrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"github.com/rtrzebinski/simple-memorizer-4/internal/probe"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/auth"
	intgrpc "github.com/rtrzebinski/simple-memorizer-4/internal/services/auth/grpc"
	"github.com/rtrzebinski/simple-memorizer-4/internal/signal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type config struct {
	ServerAddr      string        `envconfig:"SERVER_ADDRESS" default:":50051"`
	ProbeAddr       string        `envconfig:"PROBE_ADDRESS" default:"0.0.0.0:9092"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	Keycloak        struct {
		URL          string `envconfig:"KC_URL" default:"http://localhost:8180"`
		Realm        string `envconfig:"KC_REALM" default:"realm-dev"`
		ClientID     string `envconfig:"KC_CLIENT_ID" default:"client-id-dev"`
		ClientSecret string `envconfig:"KC_CLIENT_SECRET" default:"client-secret-dev"`
	}
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
			slog.Info("version", "version", version, "service", "auth")
		} else {
			slog.Info("version unknown", "service", "auth")
		}
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
	// Start auth gRPC server
	// =========================================

	listener, err := net.Listen("tcp", cfg.ServerAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	healthServer := health.NewServer()

	// register health check service
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	// register auth service
	kc := gocloak.NewClient(cfg.Keycloak.URL)
	server := intgrpc.NewServer(auth.NewService(kc, auth.Config{
		Realm:        cfg.Keycloak.Realm,
		ClientID:     cfg.Keycloak.ClientID,
		ClientSecret: cfg.Keycloak.ClientSecret,
	}))
	gengrpc.RegisterAuthServiceServer(grpcServer, server)

	healthServer.SetServingStatus("sm4-auth", grpc_health_v1.HealthCheckResponse_SERVING)

	go func() {
		slog.Info("initializing gRPC server", "addr", cfg.ServerAddr, "service", "auth")
		serverErrors <- grpcServer.Serve(listener)
	}()

	// =========================================
	// Start probes
	// =========================================

	probeServer := probe.SetupProbeServer(
		cfg.ProbeAddr,
		probe.KeycloakChecker(kc, cfg.Keycloak.Realm, cfg.Keycloak.ClientID, cfg.Keycloak.ClientSecret),
	)

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
