package http

import (
	"fmt"
	"net/http"
	"path/filepath"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

const jsonNotFoundResponse = `{"error":{"code": "NOT_FOUND","message":"resource not found"}}`
const jsonMethodNotAllowedResponse = `{"error":{"code": "METHOD_NOT_ALLOWED","message":"method not allowed"}}`

// routes sets up all the HTTP routes for the server.
// Handlers are registered from the server's handler list and loaded dynamically.
func (s *Server) routes() {
	if s.oTelEnabled {
		// Instrument the mux router for distributed traces
		s.router.Use(otelmux.Middleware("fightpicker-server"))
	}

	for _, mw := range s.middlewares {
		s.router.Use(mw)
	}

	// Register routes from all handlers
	for _, handler := range s.handlers {
		handler.RegisterRoutes(s.router)
	}

	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			s.logger.ErrorContext(r.Context(), "failed to write health response", "error", err)
		}
	})

	s.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, jsonNotFoundResponse, http.StatusNotFound)
	})

	s.router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, jsonMethodNotAllowedResponse, http.StatusMethodNotAllowed)
	})

	// static route for default profile picture when user has not set one
	s.router.HandleFunc(
		"/static/users/default-profile-picture",
		func(w http.ResponseWriter, r *http.Request) {
			// get the gender query param
			gender := r.URL.Query().Get("gender")
			if gender == "" {
				gender = "male"
			}

			// set content type and cache control headers
			w.Header().Set("Content-Type", "image/webp")
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

			// statically serve the appropriate default prof pic
			http.ServeFile(
				w,
				r,
				filepath.Join("img", "profile_pictures", fmt.Sprintf("default-profile-picture-%s.webp", gender)),
			)
		},
	).Name("static.users.get.default-profile-picture").Methods(http.MethodGet)
}
