package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/OJOMB/fightpicker/internal/http/apiresponder"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/id"
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
	*v1.Handler
	service    Service
	pathPrefix string
}

// New creates a new Handler for authentication-related endpoints.
func New(service Service, idTool id.UUID7Parser, ctxTool contextual.ContextProvider, logger logs.Logger) (*Handler, error) {
	if logger == nil {
		return nil, v1.ErrLoggerIsNil
	}

	if idTool == nil {
		return nil, v1.ErrIDToolIsNil
	}

	if ctxTool == nil {
		return nil, v1.ErrContextToolIsNil
	}

	responder := apiresponder.NewJSONResponder(ctxTool, classifyError, logger.With("component", "handler_users_v1"))
	return &Handler{
		Handler:    v1.NewHandler(idTool, responder),
		service:    service,
		pathPrefix: pathPrefix,
	}, nil
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
func (h *Handler) RegisterRoutes(m *mux.Router) {
	// POST /api/v1/auth/login - log in a user
	h.AddRoute(m, fmt.Sprintf("%s/login", h.pathPrefix), http.MethodPost, v1.EndpointNameV1AuthLogin, h.login(h.service))

	// POST /api/v1/auth/refresh - refresh authentication tokens
	h.AddRoute(m, fmt.Sprintf("%s/refresh", h.pathPrefix), http.MethodPost, v1.EndpointNameV1AuthRefresh, h.refresh(h.service))

	// POST /api/v1/auth/logout - log out a user
	h.AddRoute(m, fmt.Sprintf("%s/logout", h.pathPrefix), http.MethodPost, v1.EndpointNameV1AuthLogout, h.logout(h.service))
}
