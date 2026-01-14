package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/pkg/errors"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
)

type UserLister interface {
	ListUsers(ctx context.Context, pageSize int, lastSeenID *uuid.UUID) ([]service.User, int, error)
}

type UserGetterByEmail interface {
	GetUserByEmail(ctx context.Context, email string) (service.User, error)
}

type UserGetterByUsername interface {
	GetUserByUsername(ctx context.Context, username string) (service.User, error)
}

type UserSearcher interface {
	UserLister
	UserGetterByEmail
	UserGetterByUsername
}

// listUsers handles the HTTP GET request for the v1 list_users endpoint.
func (h *Handler) listUsers(svc UserSearcher) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		// parse query parameters for pagination
		query := r.URL.Query()

		email := query.Get(v1.QueryParamEmail)
		username := query.Get(v1.QueryParamUsername)

		if email != "" || username != "" {
			if email != "" && username != "" {
				return errors.Wrap(
					v1.ErrIncompatibleParameters,
					"you can search by either email or username, not both",
				)
			}

			var user service.User
			var err error
			if email != "" {
				user, err = svc.GetUserByEmail(ctx, email)
			} else {
				user, err = svc.GetUserByUsername(ctx, username)
			}
			if err != nil {
				return err
			}

			resp := dtos.ListUsersResponse{
				Users:    []dtos.UserResponse{userIDOToDTO(user)},
				PageSize: 1,
			}

			h.WriteJSON(ctx, w, http.StatusOK, resp)
			return nil
		}

		pageSize, lastSeenID, err := h.ParsePaginationParams(r)
		if err != nil {
			return err
		}

		users, totalCount, err := svc.ListUsers(ctx, pageSize, lastSeenID)
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
