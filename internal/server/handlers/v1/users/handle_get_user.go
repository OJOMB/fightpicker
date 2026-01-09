package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"

	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UserGetter interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (service.User, error)
}

// getUser handles the HTTP POST request for the v1 create_user endpoint.
func (h *Handler) getUser(svc UserGetter, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersGet)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := h.parseUserID(r)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		user, err := svc.GetUserByID(ctx, userID)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		h.writeJSON(ctx, w, logger, http.StatusOK, userIDOToDTO(user))
	}
}
