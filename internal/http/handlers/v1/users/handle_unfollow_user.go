package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
)

type UserUnfollower interface {
	UnfollowUser(ctx context.Context, followeeID uuid.UUID) error
}

// unfollowUser handles the HTTP DELETE request for the v1 unfollow_user endpoint.
func (h *Handler) unfollowUser(svc UserUnfollower) v1.AppHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		followeeID, err := h.ParseID(r, v1.QueryParamUserID)
		if err != nil {
			return err
		}

		if err := svc.UnfollowUser(ctx, followeeID); err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}
