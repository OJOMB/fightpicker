package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

// RotateRefreshTokens revokes the old refresh token identified by oldTokenHash and
// creates a new refresh token with the provided parameters in a single transaction.
func (r *Repo) RotateRefreshTokens(ctx context.Context, oldTokenHash, newTokenHash string, newJTI, userID id.UUID7, ipAddress, userAgent string, newExpiresAt time.Time) error {
	if newJTI == id.UUID7Nil {
		return ErrInvalidJTI
	}

	if userID == id.UUID7Nil {
		return ErrInvalidUserID
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	// rollback is a no-op if the transaction is already committed so this is safe
	defer tx.Rollback(ctx)

	qs := r.client.WithTx(tx)

	newID := r.id.Generate()

	// store the new refresh token
	params := postgres.StoreRefreshTokenParams{
		ID:        newID,
		UserID:    userID,
		TokenHash: newTokenHash,
		Jti:       newJTI,
		ExpiresAt: pgtype.Timestamptz{
			Time:  newExpiresAt.UTC(),
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
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}

	if err := qs.StoreRefreshToken(ctx, params); err != nil {
		return dbErrorToServiceError(err)
	}

	// update the old refresh token to be revoked
	if err := qs.RevokeAndRotateRefreshTokenByHash(ctx, postgres.RevokeAndRotateRefreshTokenByHashParams{
		TokenHash: oldTokenHash,
		ReplacedBy: pgtype.UUID{
			Bytes: newID,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}); err != nil {
		return dbErrorToServiceError(err)
	}

	return tx.Commit(ctx)
}
