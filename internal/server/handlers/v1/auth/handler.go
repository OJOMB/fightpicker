package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

const pathPrefix = "/api/v1/auth"

// Service defines the interface for authentication-related operations.
type Service interface {
	UserLoginner
	UserAuthenticationRefresher
	UserLogouter
}

// Handler handles authentication-related HTTP requests.
type Handler struct {
	service    Service
	pathPrefix string
}

// New creates a new Handler for authentication-related endpoints.
func New(service Service) *Handler {
	return &Handler{
		service:    service,
		pathPrefix: pathPrefix,
	}
}

// RegisterRoutes registers the authentication-related routes with the given mux router.
func (h *Handler) RegisterRoutes(mux *mux.Router, logger logs.Logger) {
	mux.Handle(
		fmt.Sprintf("%s/login", h.pathPrefix),
		otelhttp.NewHandler(
			h.login(h.service, logger),
			v1.EndpointNameV1AuthLogin,
		),
	).Name(v1.EndpointNameV1AuthLogin).
		Methods(http.MethodPost)

	mux.Handle(
		fmt.Sprintf("%s/refresh", h.pathPrefix),
		otelhttp.NewHandler(
			h.refresh(h.service, logger),
			v1.EndpointNameV1AuthRefresh,
		),
	).Name(v1.EndpointNameV1AuthRefresh).
		Methods(http.MethodPost)

	mux.Handle(
		fmt.Sprintf("%s/logout", h.pathPrefix),
		otelhttp.NewHandler(
			h.logout(h.service, logger),
			v1.EndpointNameV1AuthLogout,
		),
	).Name(v1.EndpointNameV1AuthLogout).
		Methods(http.MethodPost)
}
