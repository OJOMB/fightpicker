package users

import (
	"context"

	"github.com/gofrs/uuid"

	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

// ListFollowers retrieves a paginated list of followers for the given user in descending order (newest first).
// it also returns the total count of followers for that user.
func (r *Repo) ListFollowers(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]service.User, int, error) {
	if lastSeenID == nil {
		// Use uuidSentinelMax to represent the starting point when no lastSeenID is provided
		// this means we don't need a special case in the query for the first page
		// we're using UUID7 which sorts lexicographically, so uuidSentinelMax is greater than any valid UUID7
		lastSeenID = &uuidSentinelMax
	}

	count, err := r.dbClient.CountFollowers(ctx, userID)
	if err != nil {
		return nil, 0, dbErrorToServiceError(err)
	} else if count == 0 {
		// user not followed by anyone return empty list
		return []service.User{}, 0, nil
	}

	args := postgres.ListFollowersParams{
		FolloweeID: userID,
		Limit:      int32(pageSize),
		ID:         *lastSeenID,
	}

	rows, err := r.dbClient.ListFollowers(ctx, args)
	if err != nil {
		return nil, 0, dbErrorToServiceError(err)
	}

	users := make([]service.User, len(rows))
	for i, row := range rows {
		users[i] = listFollowersRowDBOToUserIDO(row)
	}

	return users, int(count), nil
}
