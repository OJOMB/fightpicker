package users

import (
	"context"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/id"
)

type UserFollowerLister interface {
	ListFollowers(ctx context.Context, userID id.UUID7, pageSize int, lastSeenID *id.UUID7) ([]service.User, int, error)
}

// listFollowers handles the HTTP GET request for the v1 list_followers endpoint.
func (h *Handler) listFollowers(svc UserFollowerLister) v1.HandlerFunc {
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

		users, totalCount, err := svc.ListFollowers(ctx, userID, pageSize, lastSeenID)
		if err != nil {
			return err
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

		h.WriteJSON(ctx, w, http.StatusOK, resp)

		return nil
	}
}
