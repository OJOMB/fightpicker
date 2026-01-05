package users

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/OJOMB/fightpicker/internal/server/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UserFolloweeLister interface {
	ListFollowees(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]service.User, error)
}

// listFollowees handles the HTTP GET request for the v1 list_followees endpoint.
func (h *Handler) listFollowees(svc UserFolloweeLister, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersListFollowees)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// parse query parameters for pagination
		query := r.URL.Query()
		pageSizeStr := query.Get(v1.QueryParamPageSize)
		lastSeenIDStr := query.Get(v1.QueryParamLastSeenID)
		userIDStr := query.Get(v1.QueryParamUserID)
		userID, err := uuid.FromString(userIDStr)
		if err != nil {
			logger.ErrorContext(ctx, "invalid user_id", "error", err)
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		var lastSeenID *uuid.UUID
		if lastSeenIDStr != "" {
			var err error
			lsID, err := uuid.FromString(lastSeenIDStr)
			if err != nil {
				logger.ErrorContext(ctx, "invalid last_seen_id", "error", err)
				http.Error(w, "invalid last_seen_id", http.StatusBadRequest)
				return
			}

			lastSeenID = &lsID
		}

		var pageSize int
		if pageSizeStr == "" {
			pageSize = v1.DefaultPageSize
		} else {
			pageSize, err = strconv.Atoi(pageSizeStr)
			if err != nil || pageSize < 0 {
				logger.ErrorContext(ctx, "invalid page_size", "error", err)
				http.Error(w, "invalid page_size", http.StatusBadRequest)
				return
			}

			if pageSize > v1.MaxPageSize {
				pageSize = v1.MaxPageSize
			}
		}

		users, err := svc.ListFollowees(ctx, userID, pageSize, lastSeenID)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		resp := dtos.ListUsersResponse{
			Users:    make([]dtos.UserResponse, len(users)),
			PageSize: len(users),
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
