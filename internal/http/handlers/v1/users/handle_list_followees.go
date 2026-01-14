package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
)

type UserFolloweeLister interface {
	ListFollowees(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]service.User, int, error)
}

// listFollowees handles the HTTP GET request for the v1 list_followees endpoint.
func (h *Handler) listFollowees(svc UserFolloweeLister) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		userID, err := h.ParseID(r, v1.QueryParamUserID)
		if err != nil {
			return err
		}

		pageSize, lastSeenID, err := h.ParsePaginationParams(r)
		if err != nil {
			return err
		}

		users, totalCount, err := svc.ListFollowees(ctx, userID, pageSize, lastSeenID)
		if err != nil {
			return err
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

		h.WriteJSON(ctx, w, http.StatusOK, resp)

		return nil
	}
}
