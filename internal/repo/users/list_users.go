package users

import (
	"context"

	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

func (r *Repo) ListUsers(ctx context.Context, pageSize int, lastSeenID *id.UUID7) ([]service.User, int, error) {
	if lastSeenID == nil {
		// Use id.UUID7Nil to represent the starting point when no lastSeenID is provided
		// this means we don't need a special case in the query for the first page
		// we're using UUID7 which sorts lexicographically, so id.UUID7Nil is less than any valid UUID7
		zeroUUID := id.UUID7Nil
		lastSeenID = &zeroUUID
	}

	// Get total count of users
	totalCountInt64, err := r.dbClient.CountUsers(ctx)
	if err != nil {
		return nil, 0, dbErrorToServiceError(err)
	} else if totalCountInt64 < 1 {
		return []service.User{}, 0, nil
	}

	args := postgres.ListUsersParams{
		ID:    *lastSeenID,
		Limit: int32(pageSize),
	}

	rows, err := r.dbClient.ListUsers(ctx, args)
	if err != nil {
		return nil, 0, dbErrorToServiceError(err)
	}

	users := make([]service.User, len(rows))
	for i, row := range rows {
		users[i] = listUsersRowDBOToUserIDO(row)
	}

	return users, int(totalCountInt64), nil
}
