package users

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

// FollowUser creates a follow relationship between the follower and followee in the database.
func (r *Repo) FollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error {
	params := postgres.FollowUserParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
		CreatedAt: pgtype.Timestamptz{
			Time:  r.dateTimeTool.Now().UTC(),
			Valid: true,
		},
		CreatedBy: pgtype.UUID{
			Bytes: followerID,
			Valid: true,
		},
	}

	if err := r.dbClient.FollowUser(ctx, params); err != nil {
		// if the follow relationship already exists, we can ignore the error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "followers_pkey" {
				return nil
			}
		}

		return dbErrorToServiceError(err)
	}

	return nil
}
