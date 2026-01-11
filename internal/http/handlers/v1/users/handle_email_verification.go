package users

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type EmailVerifier interface {
	VerifyEmailByToken(ctx context.Context, token string) error
}

// verifyEmail handles the HTTP POST request for the v1 create_user endpoint.
func (h *Handler) verifyEmail(svc EmailVerifier, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersEmailVerification)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// extract token from query parameters
		token := r.URL.Query().Get(v1.QueryParamEmailVerificationToken)
		if token == "" {
			h.writeError(ctx, w, errors.Wrap(v1.ErrMissingRequiredQueryParameter, "token"), logger)
			return
		}

		if err := svc.VerifyEmailByToken(ctx, token); err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
