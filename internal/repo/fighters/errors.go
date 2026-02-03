package fighters

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

// private initialization errors
var (
	// errNilLogger is returned when a nil logger is provided to the repository during initialization.
	errNilLogger = fmt.Errorf("logger cannot be nil")
	// errNilDBPool is returned when a nil database pool is provided to the repository during initialization.
	errNilDBPool = fmt.Errorf("database pool cannot be nil")
	// errNilDBClient is returned when a nil database client is provided to the repository during initialization.
	errNilDBClient = fmt.Errorf("database client cannot be nil")
	// errNilNowFunc is returned when a nil now function is provided to the repository during initialization.
	errNilNowFunc = fmt.Errorf("now function cannot be nil")
	// ErrInternalError is returned when the repo experiences an unexpected error.
)

// public runtime errors
var (
	ErrInternalError = fmt.Errorf("internal error")
	// ErrFighterNotFound is returned when a fighter with the specified ID does not exist.
	ErrFighterNotFound = fmt.Errorf("fighter not found")
)

func dbErrorToServiceError(err error) error {
	// handle pgconn.PgError and convert to service errors as needed
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrFighterNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503": // foreign_key_violation
			// could be a missing related resource
			return ErrFighterNotFound
		}
	}

	return err
}
