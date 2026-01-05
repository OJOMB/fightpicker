package users

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

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

		id := mux.Vars(r)[v1.QueryParamUserID]
		userID, err := uuid.FromString(id)
		if err != nil {
			logger.DebugContext(ctx, "invalid user_id parameter", "error", err, "user_id", id)
			h.writeError(ctx, w, errors.Wrap(v1.ErrInvalidUUID, v1.QueryParamUserID), logger)
			return
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			logger.DebugContext(ctx, "failed to read request body", "error", err)
			h.writeError(ctx, w, v1.ErrUnreadableRequestBody, logger)
			return
		}
		defer r.Body.Close()

		var userUpdateReq dtos.UserUpdateReq
		if err := json.Unmarshal(reqBody, &userUpdateReq); err != nil {
			logger.DebugContext(ctx, "failed to parse request body", "error", err)
			h.writeError(ctx, w, v1.ErrInvalidJSONRequestBody, logger)
			return
		}

		updatedUser, err := svc.UpdateUser(ctx, userID, userUpdateDTOToIDO(userUpdateReq))
		if err != nil {
			logger.ErrorContext(ctx, "failed to update user", "error", err)
			h.writeError(ctx, w, err, logger)
			return
		}

		respBody, err := json.Marshal(userIDOToDTO(updatedUser))
		if err != nil {
			logger.ErrorContext(ctx, "failed to marshal response body", "error", err)
			http.Error(w, "failed to marshal response body", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(respBody); err != nil {
			logger.ErrorContext(ctx, "failed to write response body", "error", err)
		}
	}
}
