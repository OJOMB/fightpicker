package users

import (
	"context"
	"net/http"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/id"
)

type UserGetter interface {
	GetUserByID(ctx context.Context, id id.UUID7) (service.User, error)
}

// getUser handles the HTTP POST request for the v1 create_user endpoint.
func (h *Handler) getUser(svc UserGetter) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		userID, err := h.ParseID(r, v1.QueryParamUserID)
		if err != nil {
			return err
		}

		user, err := svc.GetUserByID(ctx, userID)
		if err != nil {
			return err
		}

		h.WriteJSON(ctx, w, http.StatusOK, userIDOToDTO(user))

		return nil
	}
}
