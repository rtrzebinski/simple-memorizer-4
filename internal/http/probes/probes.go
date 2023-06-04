package http

import (
	"database/sql"
	"net/http"
)

func SetupProbeServer(addr string, db *sql.DB) *http.Server {
	r := http.NewServeMux()

	r.HandleFunc("/healthz", Healthz)
	r.HandleFunc("/readyz", Readyz(db))

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
func Readyz(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if db.Ping() != nil {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
