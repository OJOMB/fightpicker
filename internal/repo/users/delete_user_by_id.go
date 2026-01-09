package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"
)

// DeleteUserByID deletes a user by their ID.
func (r *Repo) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	if err := r.dbClient.DeleteUserByID(ctx, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}

		r.logger.ErrorContext(ctx, "failed to delete user by id", "error", err, "user_id", userID)

		return dbErrorToServiceError(err)
	}

	return nil
}
