package http

import (
	"log/slog"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/auth"
)

// Auth returns a middleware that authenticates incoming HTTP requests using
// access and refresh tokens stored in HttpOnly cookies.
func Auth(v TokenVerifier, r TokenRefresher, secure bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ca, err := req.Cookie("access_token")

			// access_token cookie missing (expired) - refresh
			if err != nil || ca.Value == "" {
				cr, err := req.Cookie("refresh_token")

				if err != nil || cr.Value == "" {
					slog.Error("no access_token and no refresh_token cookie")
					http.Error(res, "unauthorized", http.StatusUnauthorized)

					return
				}

				t, err := r.Refresh(req.Context(), cr.Value)
				if err != nil {
					slog.Error("refresh failed", err.Error())
					http.Error(res, "unauthorized", http.StatusUnauthorized)

					// clear cookies
					http.SetCookie(res, &http.Cookie{
						Name:     "refresh_token",
						Value:    "",
						Path:     "/",
						MaxAge:   -1,
						HttpOnly: true,
					})
					http.SetCookie(res, &http.Cookie{
						Name:     "access_token",
						Value:    "",
						Path:     "/",
						MaxAge:   -1,
						HttpOnly: true,
					})

					return
				}

				http.SetCookie(res, &http.Cookie{
					Name:     "refresh_token",
					Value:    t.RefreshToken,
					Path:     "/",
					MaxAge:   t.RefreshExpiresIn,
					Secure:   secure,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				})

				http.SetCookie(res, &http.Cookie{
					Name:     "access_token",
					Value:    t.AccessToken,
					Path:     "/",
					MaxAge:   t.ExpiresIn,
					Secure:   secure,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				})

				user, err := v.VerifyAndUser(req.Context(), t.AccessToken)
				if err != nil {
					slog.Error("verify after refresh failed", err.Error())
					http.Error(res, "unauthorized", http.StatusUnauthorized)

					return
				}

				// add user to ctx to be accessed by handlers
				ctx := auth.ContextWithUser(req.Context(), user)
				next.ServeHTTP(res, req.WithContext(ctx))

				return
			}

			// access_token cookie present - verify -> if needed refresh
			user, err := v.VerifyAndUser(req.Context(), ca.Value)
			if err != nil {
				cr, rerr := req.Cookie("refresh_token")

				if rerr != nil || cr.Value == "" {
					slog.Error("verify failed and no refresh_token cookie", err.Error())
					http.Error(res, "unauthorized", http.StatusUnauthorized)

					return
				}

				t, rerr := r.Refresh(req.Context(), cr.Value)
				if rerr != nil {
					slog.Error("refresh failed", rerr.Error())
					http.Error(res, "unauthorized", http.StatusUnauthorized)

					return
				}

				http.SetCookie(res, &http.Cookie{
					Name:     "refresh_token",
					Value:    t.RefreshToken,
					Path:     "/",
					MaxAge:   t.RefreshExpiresIn,
					Secure:   secure,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				})

				http.SetCookie(res, &http.Cookie{
					Name:     "access_token",
					Value:    t.AccessToken,
					Path:     "/",
					MaxAge:   t.ExpiresIn,
					Secure:   secure,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				})

				user, err = v.VerifyAndUser(req.Context(), t.AccessToken)
				if err != nil {
					slog.Error("verify after refresh failed", err.Error())
					http.Error(res, "unauthorized", http.StatusUnauthorized)

					return
				}
			}

			// add user to ctx to be accessed by handlers
			ctx := auth.ContextWithUser(req.Context(), user)
			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
