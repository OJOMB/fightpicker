package fighters

import (
	"context"

	service "github.com/OJOMB/fightpicker/internal/service/fighters"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

func (r *Repo) ListFighters(ctx context.Context, pageSize int, lastSeenID *id.UUID7) ([]service.Fighter, int, error) {
	if lastSeenID == nil {
		// Use id.UUID7Nil to represent the starting point when no lastSeenID is provided
		// this means we don't need a special case in the query for the first page
		// we're using UUID7 which sorts lexicographically, so id.UUID7Nil is less than any valid UUID7
		zeroUUID := id.UUID7Nil
		lastSeenID = &zeroUUID
	}

	// Get total count of fighters
	totalCountInt64, err := r.dbClient.CountFighters(ctx)
	if err != nil {
		return nil, 0, dbErrorToServiceError(err)
	} else if totalCountInt64 < 1 {
		return []service.Fighter{}, 0, nil
	}

	args := postgres.ListFightersParams{
		ID:    *lastSeenID,
		Limit: int32(pageSize),
	}

	rows, err := r.dbClient.ListFighters(ctx, args)
	if err != nil {
		return nil, 0, dbErrorToServiceError(err)
	}

	fighters := make([]service.Fighter, len(rows))
	for i, row := range rows {
		fighters[i] = fighterDBOtoFighterIDO(row)
	}

	return fighters, int(totalCountInt64), nil
}
