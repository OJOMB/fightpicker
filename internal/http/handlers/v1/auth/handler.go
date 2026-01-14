package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
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
	v1.Handler
	service    Service
	pathPrefix string
}

// New creates a new Handler for authentication-related endpoints.
func New(service Service, logger logs.Logger) *Handler {
	return &Handler{
		Handler: v1.Handler{
			Logger: logger.With("component", "http_handler_v1_auth"),
		},
		service:    service,
		pathPrefix: pathPrefix,
	}
}

func (h *Handler) generateRefreshCookie(refreshToken string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     v1.CookieKeyRefreshToken,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     h.pathPrefix, // restrict where it's sent
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
	}
}

// RegisterRoutes registers the authentication-related routes with the given mux router.
func (h *Handler) RegisterRoutes(mux *mux.Router) {
	mux.Handle(
		fmt.Sprintf("%s/login", h.pathPrefix),
		otelhttp.NewHandler(
			h.ToHandler(
				h.login(h.service),
				classifyError,
			),
			v1.EndpointNameV1AuthLogin,
		),
	).Name(v1.EndpointNameV1AuthLogin).
		Methods(http.MethodPost)

	mux.Handle(
		fmt.Sprintf("%s/refresh", h.pathPrefix),
		otelhttp.NewHandler(
			h.ToHandler(
				h.refresh(h.service),
				classifyError,
			),
			v1.EndpointNameV1AuthRefresh,
		),
	).Name(v1.EndpointNameV1AuthRefresh).
		Methods(http.MethodPost)

	mux.Handle(
		fmt.Sprintf("%s/logout", h.pathPrefix),
		otelhttp.NewHandler(
			h.ToHandler(
				h.logout(h.service),
				classifyError,
			),
			v1.EndpointNameV1AuthLogout,
		),
	).Name(v1.EndpointNameV1AuthLogout).
		Methods(http.MethodPost)
}
