package users

import (
	"context"
	"net/http"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/id"
)

type UserDeleter interface {
	DeleteUserByID(ctx context.Context, id id.UUID7) error
}

// deleteUser handles the HTTP DELETE request for the v1 delete_user endpoint.
func (h *Handler) deleteUser(svc UserDeleter) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		userID, err := h.ParseID(r, v1.QueryParamUserID)
		if err != nil {
			return err
		}

		if err = svc.DeleteUserByID(ctx, userID); err != nil {
			h.WriteError(ctx, w, classifyError(err))
			return err
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}
