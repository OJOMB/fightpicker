package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"

	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/id"
)

// CreateUser converts a User IDO to a DBO and calls the repo to create the user in the database
// It also assigns the default "user" role to the newly created user.
func (r *Repo) UpdateUser(ctx context.Context, userID id.UUID7, updates service.UserUpdate) (service.User, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return service.User{}, err
	}

	// rollback is a no-op if the transaction is already committed so this is safe
	defer tx.Rollback(ctx)

	qs := r.dbClient.WithTx(tx)

	dbUser, err := qs.GetUserByID(ctx, userID)
	if err != nil {
		return service.User{}, dbErrorToServiceError(err)
	}

	user := userByIDDBOToUserIDO(dbUser)
	updatedUser := user.Update(updates)

	if user.Equals(updatedUser) {
		// no updates to be made
		return user, nil
	}

	dbUpdate := UserIDOtoUpdateUserParamsDBO(updatedUser)
	if err := qs.UpdateUserByID(ctx, dbUpdate); err != nil {
		return service.User{}, dbErrorToServiceError(err)
	}

	if updates.Username != nil {
		// check that the new username is not already taken
		user, err := qs.GetUserByUsername(ctx, *updates.Username)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return service.User{}, ErrInternalError
		}

		if user.ID != id.UUID7Nil {
			return service.User{}, ErrUsernameTaken
		}
	}

	if updates.Email != nil {
		// check that the new email is not already taken
		user, err := qs.GetUserByEmail(ctx, *updates.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return service.User{}, ErrInternalError
		}

		if user.ID != id.UUID7Nil {
			return service.User{}, ErrEmailTaken
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return service.User{}, err
	}

	return updatedUser, nil
}
