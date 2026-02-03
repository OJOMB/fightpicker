package fighters

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type CacheClient interface {
	redis.Cmdable
}

type Client interface {
	GetFighterByID(ctx context.Context, id id.UUID7) (postgres.Fighter, error)
	IngestFighters(ctx context.Context, arg postgres.IngestFightersParams) ([]postgres.IngestFightersRow, error)
	WithTx(tx pgx.Tx) *postgres.Queries
}

type Repo struct {
	pool           *pgxpool.Pool
	client         Client
	cachingEnabled bool
	cache          CacheClient
	id             id.UUID7Parser
	now            datetimes.Now
	logger         logs.Logger
}

func New(pool *pgxpool.Pool, client Client, cache CacheClient, now datetimes.Now, logger logs.Logger) (*Repo, error) {
	if logger == nil {
		return nil, errNilLogger
	}

	if pool == nil {
		return nil, errNilDBPool
	}

	if client == nil {
		return nil, errNilDBClient
	}

	if now == nil {
		return nil, errNilNowFunc
	}

	cacheingEnabled := cache != nil

	return &Repo{
		pool:           pool,
		client:         client,
		cache:          cache,
		cachingEnabled: cacheingEnabled,
		logger:         logger,
		now:            now,
	}, nil
}
