package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user service.User) (service.User, error)
}

// createUser handles the HTTP POST request for the v1 create_user endpoint.
func (h *Handler) createUser(svc UserCreator, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersCreate)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		defer r.Body.Close()
		var userCreateReq dtos.UserCreateReq
		if err := json.NewDecoder(r.Body).Decode(&userCreateReq); err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		createdUser, err := svc.CreateUser(ctx, userCreateRequestDTOtoIDO(userCreateReq))
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		h.writeJSON(ctx, w, logger, http.StatusCreated, userIDOToDTO(createdUser))
	}
}
