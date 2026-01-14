package auth

import (
	"context"
	"net/http"
	"time"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
)

// UserLogouter defines the interface for logging out users.
type UserLogouter interface {
	Logout(ctx context.Context, refreshToken string) error
}

// logout handles the logout HTTP request.
func (h *Handler) logout(svc UserLogouter) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		// get refresh token from cookie
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			return err
		}

		if err := svc.Logout(ctx, cookie.Value); err != nil {
			return err
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
		return nil
	}
}
