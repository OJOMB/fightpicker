package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/id"
)

type UserUpdater interface {
	UpdateUser(ctx context.Context, id id.UUID7, updates service.UserUpdate) (service.User, error)
}

// updateUser handles the HTTP PATCH request for the v1 update_user endpoint.
func (h *Handler) updateUser(svc UserUpdater) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		userID, err := h.ParseID(r, v1.QueryParamUserID)
		if err != nil {
			return err
		}

		var userUpdateReq dtos.UserUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&userUpdateReq); err != nil {
			return v1.ErrInvalidJSONRequestBody
		}

		updatedUser, err := svc.UpdateUser(ctx, userID, userUpdateDTOToIDO(userUpdateReq))
		if err != nil {
			return err
		}

		h.Write(ctx, w, http.StatusOK, userIDOToDTO(updatedUser))

		return nil
	}
}
