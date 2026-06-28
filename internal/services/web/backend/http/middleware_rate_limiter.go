package http

import (
	"net/http"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// rateLimit takes an http.Handler and returns a new http.Handler wrapped with rate limiting
func rateLimit() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		rate := limiter.Rate{
			Period: 1 * time.Minute,
			Limit:  20,
		}

		store := memory.NewStoreWithOptions(limiter.StoreOptions{
			CleanUpInterval: 5 * time.Minute,
		})

		// K8s Ingress must be configured to set the X-Forwarded-For header
		// https://github.com/ulule/limiter#limiter-behind-a-reverse-proxy
		instance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))

		m := stdlib.NewMiddleware(instance)
		m.OnLimitReached = func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}

		return m.Handler(next)
	}
}
