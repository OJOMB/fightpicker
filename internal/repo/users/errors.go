package users

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrUserNotFound is returned when a user is not found in the database.
	ErrUserNotFound = fmt.Errorf("user not found")
	// ErrEmailTaken is returned when a user email is already taken.
	ErrEmailTaken = fmt.Errorf("email already taken")
	// ErrUsernameTaken is returned when a user username is already taken.
	ErrUsernameTaken = fmt.Errorf("username already taken")
	// ErrDefaultRoleNotFound indicates that the rbac system has not been setup properly.
	ErrDefaultRoleNotFound = fmt.Errorf("default role not found")
	// ErrInternalError is returned when the repo experiences an unexpected error.
	ErrInternalError = fmt.Errorf("internal error")
)

func dbErrorToServiceError(err error) error {
	// handle pgconn.PgError and convert to service errors as needed
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			switch pgErr.ConstraintName {
			case "users_email_key":
				return ErrEmailTaken
			case "users_username_key":
				return ErrUsernameTaken
			}

			// should not really happen here but just in case
			if pgErr.TableName == "user_roles" && pgErr.ColumnName == "role_id" {
				return ErrDefaultRoleNotFound
			}
		}
	}

	return err
}
