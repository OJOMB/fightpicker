package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UserFollower interface {
	FollowUser(ctx context.Context, followeeID uuid.UUID) error
}

// followUser handles the HTTP PUT request for the v1 follow_user endpoint.
func (h *Handler) followUser(svc UserFollower, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersFollow)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		followeeID, err := h.parseUserID(r)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		if err := svc.FollowUser(ctx, followeeID); err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
