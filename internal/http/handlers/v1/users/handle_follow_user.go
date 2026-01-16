package users

import (
	"context"
	"net/http"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/id"
)

type UserFollower interface {
	FollowUser(ctx context.Context, followeeID id.UUID7) error
}

// followUser handles the HTTP PUT request for the v1 follow_user endpoint.
func (h *Handler) followUser(svc UserFollower) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		followeeID, err := h.ParseID(r, v1.QueryParamUserID)
		if err != nil {
			return err
		}

		if err := svc.FollowUser(ctx, followeeID); err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}
