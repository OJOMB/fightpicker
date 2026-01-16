package users

import (
	"context"
	"database/sql"
	"errors"

	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/id"
)

// GetUserByID retrieves a user by their ID.
func (r *Repo) GetUserByID(ctx context.Context, userID id.UUID7) (service.User, error) {
	dbUser, err := r.dbClient.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.User{}, ErrUserNotFound
		}

		return service.User{}, dbErrorToServiceError(err)
	}

	return userByIDDBOToUserIDO(dbUser), nil
}
