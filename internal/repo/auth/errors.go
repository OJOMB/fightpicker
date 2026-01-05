package auth

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrUserNotFound is returned when a user is not found in the database.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserEmailTaken is returned when a user email is already taken.
	ErrUserEmailTaken = errors.New("email already taken")
	// ErrUserUsernameTaken is returned when a user username is already taken.
	ErrUserUsernameTaken = errors.New("username already taken")

	// ErrInternalError is returned when an unexpected internal error occurs.
	ErrInternalError = errors.New("internal error")

	// ErrInvalidJTI is returned when a provided JTI is invalid.
	ErrInvalidJTI = errors.New("invalid JTI provided")

	// ErrInvalidUserID is returned when a provided user ID is invalid.
	ErrInvalidUserID = errors.New("invalid user ID provided")
)

func dbErrorToServiceError(err error) error {
	// handle pgconn.PgError and convert to service errors as needed
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			switch pgErr.ConstraintName {
			case "users_email_key":
				return ErrUserEmailTaken
			case "users_username_key":
				return ErrUserUsernameTaken
			}
		}
	}

	return err
}
