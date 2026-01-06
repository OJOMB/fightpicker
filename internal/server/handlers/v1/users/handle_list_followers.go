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

type UserFollowerLister interface {
	ListFollowers(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]service.User, int, error)
}

// listFollowers handles the HTTP GET request for the v1 list_followers endpoint.
func (h *Handler) listFollowers(svc UserFollowerLister, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersListFollowers)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := h.parseUserID(r)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		lastSeenID, err := h.parseLastSeenID(r)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		pageSize, err := h.parsePageSize(r)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		users, totalCount, err := svc.ListFollowers(ctx, userID, pageSize, lastSeenID)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		resp := dtos.ListUsersResponse{
			Users:      make([]dtos.UserResponse, len(users)),
			PageSize:   len(users),
			TotalCount: totalCount,
		}
		for i, user := range users {
			resp.Users[i] = userIDOToDTO(user)
		}

		// set LastSeenId for pagination if there are more users to fetch
		if len(users) > 0 && len(users) == pageSize {
			resp.LastSeenId = &resp.Users[len(resp.Users)-1].Id
		}

		respBody, err := json.Marshal(resp)
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
