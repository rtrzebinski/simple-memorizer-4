package http

import (
	"database/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net/http"
)

func SetupProbeServer(addr string, db *sql.DB, conn *grpc.ClientConn) *http.Server {
	r := http.NewServeMux()

	r.HandleFunc("/healthz", Healthz)
	r.HandleFunc("/readyz", Readyz(db, conn))

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

// Healthz is a liveness probe.
func Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readyz is a readiness probe.
func Readyz(db *sql.DB, conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// db connection check
		if db != nil && db.Ping() != nil {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		// grpc connection check
		if conn != nil {
			healthClient := grpc_health_v1.NewHealthClient(conn)
			response, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
				Service: "sm4-auth-service",
			})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
				return
			}
			if response.Status != grpc_health_v1.HealthCheckResponse_SERVING {
				http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
