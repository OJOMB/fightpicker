package fighters

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type CacheClient interface {
	redis.Cmdable
}

type Client interface {
	GetFighterByID(ctx context.Context, id uuid.UUID) (postgres.Fighter, error)
	WithTx(tx pgx.Tx) *postgres.Queries
}

type Repo struct {
	pool           *pgxpool.Pool
	client         Client
	cachingEnabled bool
	cache          CacheClient
	now            datetimes.Now
	logger         logs.Logger
}

func New(pool *pgxpool.Pool, client Client, cache CacheClient, now datetimes.Now, logger logs.Logger) (*Repo, error) {
	if logger == nil {
		return nil, ErrNilLogger
	}

	if pool == nil {
		return nil, ErrNilDBPool
	}

	if client == nil {
		return nil, ErrNilDBClient
	}

	if now == nil {
		return nil, ErrNilNowFunc
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
