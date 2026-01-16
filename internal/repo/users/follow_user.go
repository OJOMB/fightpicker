package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

// FollowUser creates a follow relationship between the follower and followee in the database.
// the followID is a unique identifier for the follow relationship used for pagination consistency and simplicity.
func (r *Repo) FollowUser(ctx context.Context, followID, followerID, followeeID id.UUID7) error {
	params := postgres.FollowUserParams{
		ID:         followID,
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
