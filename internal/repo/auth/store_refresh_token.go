package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

func (r *Repo) StoreRefreshToken(ctx context.Context, userID id.UUID7, jti id.UUID7, tokenHash, ipAddress, userAgent string, expiresAt time.Time) error {
	now := time.Now().UTC()
	params := postgres.StoreRefreshTokenParams{
		ID:        r.id.Generate(),
		FamilyID:  r.id.Generate(),
		UserID:    userID,
		TokenHash: tokenHash,
		Jti:       jti,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt.UTC(),
			Valid: true,
		},
		IpAddress: pgtype.Text{
			String: ipAddress,
			Valid:  ipAddress != "",
		},
		UserAgent: pgtype.Text{
			String: userAgent,
			Valid:  userAgent != "",
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
	}

	if err := r.client.StoreRefreshToken(ctx, params); err != nil {
		return r.dbErrorToRepoError(ctx, err)
	}

	return nil
}
