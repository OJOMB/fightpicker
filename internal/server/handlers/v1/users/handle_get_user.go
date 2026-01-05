package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

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
		userIDStr := mux.Vars(r)[v1.QueryParamUserID]
		userID, err := uuid.FromString(userIDStr)
		if err != nil {
			logger.DebugContext(ctx, "invalid user_id parameter", "error", err, "user_id", userIDStr)
			h.writeError(ctx, w, errors.Wrap(v1.ErrInvalidUUID, v1.QueryParamUserID), logger)
			return
		}

		user, err := svc.GetUserByID(ctx, userID)
		if err != nil {
			logger.ErrorContext(ctx, "failed to get user", "error", err)
			h.writeError(ctx, w, err, logger)
			return
		}

		respBody, err := json.Marshal(userIDOToDTO(user))
		if err != nil {
			logger.ErrorContext(ctx, "failed to marshal response body", "error", err)
			http.Error(w, "failed to marshal response body", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(respBody); err != nil {
			logger.ErrorContext(ctx, "failed to write response body", "error", err)
		}
	}
}
