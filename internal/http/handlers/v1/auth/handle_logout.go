package auth

import (
	"context"
	"net/http"
	"time"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// UserLogouter defines the interface for logging out users.
type UserLogouter interface {
	Logout(ctx context.Context, refreshToken string) error
}

// logout handles the logout HTTP request.
func (h *Handler) logout(svc UserLogouter, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1AuthLogout)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// get refresh token from cookie
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		refreshToken := cookie.Value

		if err := svc.Logout(ctx, refreshToken); err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		// Clear the refresh token cookie
		clearCookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			Path:     "/",                  // clear for all paths
			SameSite: http.SameSiteLaxMode, // or Strict if possible
			Expires:  time.Unix(0, 0),      // set expiration in the past
		}
		http.SetCookie(w, clearCookie)

		w.WriteHeader(http.StatusNoContent)
	}
}
