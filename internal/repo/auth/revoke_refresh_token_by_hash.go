package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

func (r *Repo) RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error {
	params := postgres.RevokeRefreshTokenByHashParams{
		TokenHash: tokenHash,
		UpdatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}
	if err := r.client.RevokeRefreshTokenByHash(ctx, params); err != nil {
		return dbErrorToRepoError(err)
	}

	return nil
}
