package http

import (
	"net/http"
	"net/url"
	"strings"
)

// CSRFDynamicHost returns a middleware that protects against CSRF attacks
// by rejecting unsafe HTTP requests (POST, PUT, PATCH, DELETE) coming from
// other domains. It compares the Origin or Referer header to the request’s
// Host value and only allows requests where they match.
func CSRFDynamicHost() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				reqHost := hostOnly(r.Host)

				origin := r.Header.Get("Origin")
				referer := r.Header.Get("Referer")

				headerHost := ""

				if origin != "" {
					if u, err := url.Parse(origin); err == nil {
						headerHost = hostOnly(u.Host)
					}
				} else if referer != "" {
					if u, err := url.Parse(referer); err == nil {
						headerHost = hostOnly(u.Host)
					}
				}

				if headerHost == "" || !strings.EqualFold(headerHost, reqHost) {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func hostOnly(h string) string {
	if h == "" {
		return ""
	}
	if strings.HasPrefix(h, "[") {
		if i := strings.Index(h, "]"); i != -1 {
			return h[:i+1]
		}
	}
	if i := strings.IndexByte(h, ':'); i != -1 {
		return h[:i]
	}
	return h
}
