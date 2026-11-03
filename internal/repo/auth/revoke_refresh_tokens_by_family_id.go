package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

func (r *Repo) RevokeRefreshTokensByFamilyID(ctx context.Context, familyID id.UUID7) error {
	params := postgres.RevokeRefreshTokensByFamilyIDParams{
		FamilyID: familyID,
		UpdatedAt: pgtype.Timestamptz{
			Time:  r.datetime.Now(),
			Valid: true,
		},
	}

	rowsAffected, err := r.client.RevokeRefreshTokensByFamilyID(ctx, params)
	if err != nil {
		return r.dbErrorToRepoError(ctx, err)
	}

	if rowsAffected == 0 {
		r.logger.DebugContext(ctx, "no rows affected attempting to revoke family", "family_id", familyID)
		return ErrRefreshTokenInvalidOrReused
	}

	return nil
}
