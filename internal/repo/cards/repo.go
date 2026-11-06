package cards

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type DBClient interface {
	CreateCard(ctx context.Context, arg postgres.CreateCardParams) error
}

type Repo struct {
	pool     *pgxpool.Pool
	dbClient DBClient
	logger   logs.Logger
}

func New(pool *pgxpool.Pool, client DBClient, logger logs.Logger) (*Repo, error) {
	if pool == nil {
		return nil, errNilDBPool
	}

	if client == nil {
		return nil, errNilDBClient
	}

	if logger == nil {
		return nil, errNilLogger
	}

	return &Repo{
		pool:     pool,
		dbClient: client,
		logger:   logger.With("component", "cards_repo"),
	}, nil
}
