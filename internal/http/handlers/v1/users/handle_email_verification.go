package users

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
)

type EmailVerifier interface {
	VerifyEmailByToken(ctx context.Context, token string) error
}

// verifyEmail handles the HTTP POST request for the v1 create_user endpoint.
func (h *Handler) verifyEmail(svc EmailVerifier) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		token := r.URL.Query().Get(v1.QueryParamEmailVerificationToken)
		if token == "" {
			return errors.Wrap(v1.ErrMissingRequiredQueryParameter, v1.QueryParamEmailVerificationToken)
		}

		if err := svc.VerifyEmailByToken(ctx, token); err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}
