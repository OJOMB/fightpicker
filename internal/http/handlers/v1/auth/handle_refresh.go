package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// UserAuthenticationRefresher defines the interface for refreshing user authentication tokens.
type UserAuthenticationRefresher interface {
	Refresh(ctx context.Context, refreshToken string) (string, string, time.Time, error)
}

// refresh handles the token refresh HTTP request.
func (h *Handler) refresh(svc UserAuthenticationRefresher, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1AuthRefresh)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		c, err := r.Cookie(v1.CookieKeyRefreshToken)
		if err != nil {
			logger.DebugContext(ctx, "refresh called with no refresh token in cookie", "error", err)
			h.writeError(ctx, w, ErrMissingRefreshToken, logger)
			return
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
			h.writeError(ctx, w, err, logger)
			return
		}

		// 🔐 Set refresh token as HttpOnly cookie
		refreshCookie := &http.Cookie{
			Name:     v1.CookieKeyRefreshToken,
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,                 // true in production (HTTPS)
			Path:     h.pathPrefix,         // restrict where it's sent
			SameSite: http.SameSiteLaxMode, // or Strict if possible
			Expires:  refreshTokenExpiresAt,
		}

		http.SetCookie(w, refreshCookie)

		resp := dtos.AuthResponse{AccessToken: accessToken}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.ErrorContext(ctx, "failed to write response body", "error", err)
		}
	}
}
