package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user service.User) (service.User, error)
}

// createUser handles the HTTP POST request for the v1 create_user endpoint.
func (h *Handler) createUser(svc UserCreator) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		defer r.Body.Close()
		var userCreateReq dtos.UserCreateReq
		if err := json.NewDecoder(r.Body).Decode(&userCreateReq); err != nil {
			return v1.ErrInvalidJSONRequestBody
		}

		createdUser, err := svc.CreateUser(ctx, userCreateRequestDTOtoIDO(userCreateReq))
		if err != nil {
			return err
		}

		h.WriteJSON(ctx, w, http.StatusCreated, userIDOToDTO(createdUser))

		return nil
	}
}
