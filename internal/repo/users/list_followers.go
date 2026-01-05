package users

import (
	"context"

	"github.com/gofrs/uuid"

	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

// uuidSentinelMax is a UUID greater than any valid UUID7.
// UUID7 structure:
// | 48 bits timestamp | 4 bits version | 12 bits subsec |
// | 2–3 bits variant  | 62 bits random |
// the version is always 7 (0111) and the variant bits are either 10 or 11,
// so the maximum possible UUID7 is ffffffff-ffff-7fff-bfff-ffffffffffff
var uuidSentinelMax = uuid.Must(uuid.FromString("ffffffff-ffff-7fff-bfff-ffffffffffff"))

func (r *Repo) ListFollowers(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]service.User, error) {
	if lastSeenID == nil {
		// Use uuidSentinelMax to represent the starting point when no lastSeenID is provided
		// this means we don't need a special case in the query for the first page
		// we're using UUID7 which sorts lexicographically, so uuidSentinelMax is greater than any valid UUID7
		lastSeenID = &uuidSentinelMax
	}

	args := postgres.ListFollowersParams{
		FolloweeID: userID,
		Limit:      int32(pageSize),
		ID:         *lastSeenID,
	}

	rows, err := r.dbClient.ListFollowers(ctx, args)
	if err != nil {
		return nil, dbErrorToServiceError(err)
	}

	users := make([]service.User, len(rows))
	for i, row := range rows {
		users[i] = listFollowersRowDBOToUserIDO(row)
	}

	return users, nil
}
