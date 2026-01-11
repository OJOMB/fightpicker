package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/OJOMB/fightpicker/internal/server/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// UserLoginner defines the interface for logging in users with their credentials.
type UserLoginner interface {
	Login(ctx context.Context, email, password string) (string, string, time.Time, error)
}

// login handles the HTTP POST request for the v1 login endpoint.
func (h *Handler) login(svc UserLoginner, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1AuthLogin)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// get credentials from request body
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			logger.ErrorContext(ctx, "failed to read request body", "error", err)
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var credentials dtos.LoginRequest
		if err := json.Unmarshal(reqBody, &credentials); err != nil {
			logger.ErrorContext(ctx, "failed to parse request body", "error", err)
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}

		// authenticate user
		accessToken, refreshToken, refreshTokenExpiresAt, err := svc.Login(ctx, string(credentials.Email), credentials.Password)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		// 🔐 Set refresh token as HttpOnly cookie
		refreshCookie := &http.Cookie{
			Name:     v1.CookieKeyRefreshToken,
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     h.pathPrefix,
			SameSite: http.SameSiteLaxMode,
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
