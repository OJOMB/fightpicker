package users

import (
	"context"

	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

// ListFollowees retrieves a paginated list of users that the given user follows in descending order (newest first).
// it also returns the total count of followers for that user.
func (r *Repo) ListFollowees(ctx context.Context, userID id.UUID7, pageSize int, lastSeenID *id.UUID7) ([]service.User, int, error) {
	if lastSeenID == nil {
		// Use uuidSentinelMax to represent the starting point when no lastSeenID is provided
		// this means we don't need a special case in the query for the first page
		// we're using UUID7 which sorts lexicographically, so uuidSentinelMax is greater than any valid UUID7
		lastSeenID = &id.UUID7SentinelMax
	}

	count, err := r.dbClient.CountFollowees(ctx, userID)
	if err != nil {
		return nil, 0, dbErrorToServiceError(err)
	} else if count == 0 {
		// user not following anyone return empty list
		return []service.User{}, 0, nil
	}

	args := postgres.ListFolloweesParams{
		FollowerID: userID,
		Limit:      int32(pageSize),
		ID:         *lastSeenID,
	}

	rows, err := r.dbClient.ListFollowees(ctx, args)
	if err != nil {
		return nil, 0, dbErrorToServiceError(err)
	}

	users := make([]service.User, len(rows))
	for i, row := range rows {
		users[i] = listFolloweesRowDBOToUserIDO(row)
	}

	return users, int(count), nil
}
