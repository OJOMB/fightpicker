package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
)

// UserAuthenticationRefresher defines the interface for refreshing user authentication tokens.
type UserAuthenticationRefresher interface {
	Refresh(ctx context.Context, refreshToken string) (string, string, time.Time, error)
}

// refresh handles the token refresh HTTP request.
func (h *Handler) refresh(svc UserAuthenticationRefresher) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		c, err := r.Cookie(v1.CookieKeyRefreshToken)
		if err != nil {
			return ErrMissingRefreshToken
		}

		accessToken, refreshToken, refreshTokenExpiresAt, err := svc.Refresh(ctx, c.Value)
		if err != nil {
			// if refresh fails, clear the cookie
			http.SetCookie(w, &http.Cookie{
				Name:     v1.CookieKeyRefreshToken,
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
			return err
		}

		http.SetCookie(w, h.generateRefreshCookie(refreshToken, refreshTokenExpiresAt))

		h.WriteJSON(ctx, w, http.StatusOK, dtos.AuthResponse{AccessToken: accessToken})

		return nil
	}
}
