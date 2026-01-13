package users

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

// UpdateLastLoginAtByUserID updates the last login timestamp for a user in the database.
func (r *Repo) UpdateLastLoginAtByUserID(ctx context.Context, userID uuid.UUID) error {
	params := postgres.UpdateLastLoginAtByUserIDParams{
		ID:          userID,
		LastLoginAt: pgtype.Timestamptz{Time: r.dateTimeTool.Now(), Valid: true},
	}

	if err := r.dbClient.UpdateLastLoginAtByUserID(ctx, params); err != nil {
		return dbErrorToServiceError(err)
	}

	return nil
}
