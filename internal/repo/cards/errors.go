package cards

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
)

var (
	ErrDuplicateResource = fmt.Errorf("duplicate resource")
	ErrResourceNotFound  = fmt.Errorf("resource not found")
)

func dbErrorToServiceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrResourceNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			switch pgErr.ConstraintName {
			case "cards_name_promotion_key":
				return ErrDuplicateResource
			}
		}
	}

	return err
}
