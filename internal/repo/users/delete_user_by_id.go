package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/OJOMB/fightpicker/pkg/id"
)

// DeleteUserByID deletes a user by their ID.
func (r *Repo) DeleteUserByID(ctx context.Context, userID id.UUID7) error {
	if err := r.dbClient.DeleteUserByID(ctx, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		r.logger.ErrorContext(ctx, "failed to delete user by id", "error", err, "user_id", userID)

		return dbErrorToServiceError(err)
	}

	return nil
}
