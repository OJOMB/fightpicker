package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/OJOMB/fightpicker/internal/server/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UserUpdater interface {
	UpdateUser(ctx context.Context, id uuid.UUID, updates service.UserUpdate) (service.User, error)
}

// updateUser handles the HTTP PATCH request for the v1 update_user endpoint.
func (h *Handler) updateUser(svc UserUpdater, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersUpdate)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := h.parseUserID(r)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		var userUpdateReq dtos.UserUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&userUpdateReq); err != nil {
			h.writeError(ctx, w, v1.ErrInvalidJSONRequestBody, logger)
			return
		}

		updatedUser, err := svc.UpdateUser(ctx, userID, userUpdateDTOToIDO(userUpdateReq))
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		h.writeJSON(ctx, w, logger, http.StatusOK, userIDOToDTO(updatedUser))
	}
}
