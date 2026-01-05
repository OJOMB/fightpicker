package users

import (
	"context"
	"database/sql"
	"errors"

	service "github.com/OJOMB/fightpicker/internal/service/users"
)

// GetUserByEmail retrieves a user by their email address.
func (r *Repo) GetUserByEmail(ctx context.Context, email string) (service.User, error) {
	dbUser, err := r.dbClient.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.User{}, ErrUserNotFound
		}

		return service.User{}, dbErrorToServiceError(err)
	}

	return userByEmailDBOToUserIDO(dbUser), nil
}
