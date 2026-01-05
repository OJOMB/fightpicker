package users

import (
	"context"
	"database/sql"
	"errors"

	service "github.com/OJOMB/fightpicker/internal/service/users"
)

// GetUserByUsername retrieves a user by their username.
func (r *Repo) GetUserByUsername(ctx context.Context, username string) (service.User, error) {
	dbUser, err := r.dbClient.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.User{}, ErrUserNotFound
		}

		return service.User{}, dbErrorToServiceError(err)
	}

	return userByUsernameDBOToUserIDO(dbUser), nil
}
