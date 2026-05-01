package http

import (
	"log"
	"net/http"
)

type HandlerAuthLogout struct {
	s Service
}

func NewHandlerAuthLogout(s Service) *HandlerAuthLogout {
	return &HandlerAuthLogout{
		s: s,
	}
}

func (h *HandlerAuthLogout) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	c, err := req.Cookie("refresh_token")

	if err != nil || c.Value == "" {
		log.Printf("logout: no refresh cookie (%v)", err)
		res.WriteHeader(http.StatusOK)
		return
	}

	err = h.s.Revoke(req.Context(), c.Value)
	if err != nil {
		log.Printf("logout: revoke failed: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	res.WriteHeader(http.StatusOK)
}
