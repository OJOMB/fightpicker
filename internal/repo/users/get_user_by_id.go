package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"

	service "github.com/OJOMB/fightpicker/internal/service/users"
)

// GetUserByID retrieves a user by their ID.
func (r *Repo) GetUserByID(ctx context.Context, userID uuid.UUID) (service.User, error) {
	dbUser, err := r.dbClient.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.User{}, ErrUserNotFound
		}

		return service.User{}, dbErrorToServiceError(err)
	}

	return userByIDDBOToUserIDO(dbUser), nil
}
