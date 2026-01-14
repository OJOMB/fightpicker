package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
)

// UserLoginner defines the interface for logging in users with their credentials.
type UserLoginner interface {
	Login(ctx context.Context, email, password string) (string, string, time.Time, error)
}

// login handles the HTTP POST request for the v1 login endpoint.
func (h *Handler) login(svc UserLoginner) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		var credentials dtos.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			return v1.ErrInvalidJSONRequestBody
		}

		accessToken, refreshToken, refreshTokenExpiresAt, err := svc.Login(ctx, string(credentials.Email), credentials.Password)
		if err != nil {
			return err
		}

		http.SetCookie(w, h.generateRefreshCookie(refreshToken, refreshTokenExpiresAt))

		h.WriteJSON(ctx, w, http.StatusOK, dtos.AuthResponse{AccessToken: accessToken})

		return nil
	}
}
