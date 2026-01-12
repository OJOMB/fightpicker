package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UserDeleter interface {
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
}

// deleteUser handles the HTTP DELETE request for the v1 delete_user endpoint.
func (h *Handler) deleteUser(svc UserDeleter, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersDelete)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := h.parseUserID(r)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		if err = svc.DeleteUserByID(ctx, userID); err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
