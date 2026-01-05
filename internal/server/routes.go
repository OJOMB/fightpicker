package server

import (
	"fmt"
	"net/http"
	"path/filepath"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"github.com/OJOMB/fightpicker/internal/server/middlewares"
)

// routes sets up all the HTTP routes for the server.
// Handlers are registered from the server's handler list and loaded dynamically.
func (svr *Server) routes() {
	if svr.oTelEnabled {
		// Instrument the mux router for distributed traces
		svr.router.Use(otelmux.Middleware("fightpicker-server"))
	}

	// add request logger middleware
	logRespBody := svr.env == "development" || svr.env == "local"
	svr.router.Use(middlewares.NewRequestResponseLogger(svr.logger, logRespBody, svr.oTelEnabled).Middleware)
	svr.router.Use(middlewares.NewAuthPermissionsChecker(svr.secretKey, svr.jwtValidator, svr.logger).Middleware)
	svr.router.Use(middlewares.NewContextLoader(svr.logger).Middleware)
	svr.router.Use(middlewares.NewPyroProfiler(map[string]string{"component": "server"}).Middleware)

	// Register routes from all handlers
	for _, handler := range svr.handlers {
		handler.RegisterRoutes(svr.router, svr.logger)
	}

	svr.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			svr.logger.ErrorContext(r.Context(), "failed to write health response", "error", err)
		}
	})

	svr.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
	})

	svr.router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	})

	// static route for default profile picture when user has not set one
	svr.router.HandleFunc(
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
