package probe

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/Nerzal/gocloak/v13"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const timeout = 2 * time.Second

type Checker func(ctx context.Context) error

func Healthz(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }

func Readyz(checkers []Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(checkers) == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		for _, chk := range checkers {
			err := chk(ctx)
			if err != nil {
				slog.Error("readiness check failed", slog.Any("checker", chk), slog.Any("error", err))
				http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func SetupProbeServer(addr string, checkers ...Checker) *http.Server {
	r := http.NewServeMux()
	r.HandleFunc("/healthz", Healthz)
	r.HandleFunc("/readyz", Readyz(checkers))
	return &http.Server{Addr: addr, Handler: r}
}

func DBChecker(db *sql.DB) Checker {
	return func(ctx context.Context) error {
		if db == nil {
			return errors.New("nil db")
		}
		return db.PingContext(ctx)
	}
}

func PubSubChecker(projectID string) Checker {
	return func(ctx context.Context) error {
		if projectID == "" {
			return errors.New("empty project id")
		}
		c, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			return err
		}
		return c.Close()
	}
}

func GrpcChecker(addr string) Checker {
	return func(ctx context.Context) error {
		conn, err := grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}
		defer conn.Close()
		client := healthpb.NewHealthClient(conn)
		_, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: ""})
		return err
	}
}

func KeycloakChecker(kc *gocloak.GoCloak, realm, clientID, clientSecret string) Checker {
	return func(ctx context.Context) error {
		if kc == nil {
			return errors.New("nil keycloak client")
		}

		if realm == "" {
			return errors.New("empty realm")
		}

		if clientID == "" {
			return errors.New("empty client id")
		}

		if clientSecret == "" {
			return errors.New("empty client secret")
		}

		_, err := kc.LoginClient(ctx, clientID, clientSecret, realm)

		return err
	}
}
