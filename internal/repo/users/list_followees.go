package users

import (
	"context"

	"github.com/gofrs/uuid"

	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

func (r *Repo) ListFollowees(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]service.User, error) {
	if lastSeenID == nil {
		// Use uuidSentinelMax to represent the starting point when no lastSeenID is provided
		// this means we don't need a special case in the query for the first page
		// we're using UUID7 which sorts lexicographically, so uuidSentinelMax is greater than any valid UUID7
		lastSeenID = &uuidSentinelMax
	}

	args := postgres.ListFolloweesParams{
		FollowerID: userID,
		Limit:      int32(pageSize),
		ID:         *lastSeenID,
	}

	rows, err := r.dbClient.ListFollowees(ctx, args)
	if err != nil {
		return nil, dbErrorToServiceError(err)
	}

	users := make([]service.User, len(rows))
	for i, row := range rows {
		users[i] = listFolloweesRowDBOToUserIDO(row)
	}

	return users, nil
}
