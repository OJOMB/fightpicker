package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid/v5"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

// UnfollowUser deletes a follow relationship between the follower and followee in the database.
func (r *Repo) UnfollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error {
	params := postgres.UnfollowUserParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}

	if err := r.dbClient.UnfollowUser(ctx, params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return dbErrorToServiceError(err)
	}

	return nil
}
