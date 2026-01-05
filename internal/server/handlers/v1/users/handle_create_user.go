package users

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/server/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
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

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			logger.DebugContext(ctx, "failed to read request body", "error", err)
			h.writeError(ctx, w, err, logger)
			return
		}
		defer r.Body.Close()

		var userCreateReq dtos.UserCreateReq
		if err := json.Unmarshal(reqBody, &userCreateReq); err != nil {
			logger.DebugContext(ctx, "failed to parse request body", "error", err)
			h.writeError(ctx, w, err, logger)
			return
		}

		svcUser := userCreateRequestDTOtoIDO(userCreateReq)

		createdUser, err := svc.CreateUser(ctx, svcUser)
		if err != nil {
			logger.ErrorContext(ctx, "failed to create user", "error", err)
			h.writeError(ctx, w, err, logger)
			return
		}

		respBody, err := json.Marshal(userIDOToDTO(createdUser))
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
