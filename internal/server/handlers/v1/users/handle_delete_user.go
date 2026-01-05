package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
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
		userIDStr := mux.Vars(r)[v1.QueryParamUserID]
		userID, err := uuid.FromString(userIDStr)
		if err != nil {
			logger.DebugContext(ctx, "invalid user_id parameter", "error", err, "user_id", userIDStr)
			h.writeError(ctx, w, errors.Wrap(v1.ErrInvalidUUID, v1.QueryParamUserID), logger)
			return
		}

		if err = svc.DeleteUserByID(ctx, userID); err != nil {
			logger.ErrorContext(ctx, "failed to delete user", "error", err)
			h.writeError(ctx, w, err, logger)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
