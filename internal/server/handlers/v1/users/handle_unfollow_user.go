package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UserUnfollower interface {
	UnfollowUser(ctx context.Context, followeeID uuid.UUID) error
}

// unfollowUser handles the HTTP DELETE request for the v1 unfollow_user endpoint.
func (h *Handler) unfollowUser(svc UserUnfollower, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersUnfollow)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := mux.Vars(r)[v1.QueryParamUserID]
		followeeID, err := uuid.FromString(id)
		if err != nil {
			logger.DebugContext(ctx, "invalid user_id parameter", "error", err, "followee_id", id)
			h.writeError(ctx, w, err, logger)
			return
		}

		if err := svc.UnfollowUser(ctx, followeeID); err != nil {
			logger.ErrorContext(ctx, "failed to unfollow user", "error", err)
			h.writeError(ctx, w, err, logger)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
