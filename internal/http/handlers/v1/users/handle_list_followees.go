package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UserFolloweeLister interface {
	ListFollowees(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]service.User, int, error)
}

// listFollowees handles the HTTP GET request for the v1 list_followees endpoint.
func (h *Handler) listFollowees(svc UserFolloweeLister, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersListFollowees)

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

		users, totalCount, err := svc.ListFollowees(ctx, userID, pageSize, lastSeenID)
		if err != nil {
			h.writeError(ctx, w, err, logger)
			return
		}

		resp := dtos.ListUsersResponse{
			TotalCount: totalCount,
			Users:      make([]dtos.UserResponse, len(users)),
			PageSize:   len(users),
		}

		for i, user := range users {
			resp.Users[i] = userIDOToDTO(user)
		}

		// set LastSeenId for pagination if there are more users to fetch
		if len(users) > 0 && len(users) == pageSize {
			resp.LastSeenId = &resp.Users[len(resp.Users)-1].Id
		}

		h.writeJSON(ctx, w, logger, http.StatusOK, resp)
	}
}
