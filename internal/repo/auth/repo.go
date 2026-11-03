package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type DBClient interface {
	GetUserRolesByID(ctx context.Context, userID id.UUID7) ([]string, error)
	GetUserPermissionsByID(ctx context.Context, id id.UUID7) ([]postgres.GetUserPermissionsByIDRow, error)
	StoreRefreshToken(ctx context.Context, arg postgres.StoreRefreshTokenParams) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (postgres.RefreshToken, error)
	RevokeRefreshTokenByHash(ctx context.Context, arg postgres.RevokeRefreshTokenByHashParams) error
	RevokeAndRotateRefreshTokenByHash(ctx context.Context, arg postgres.RevokeAndRotateRefreshTokenByHashParams) (int64, error)
	RevokeRefreshTokensByFamilyID(ctx context.Context, arg postgres.RevokeRefreshTokensByFamilyIDParams) (int64, error)
	WithTx(tx pgx.Tx) *postgres.Queries
}

type Repo struct {
	pool              *pgxpool.Pool
	client            DBClient
	logger            logs.Logger
	id                id.UUID7GeneratorParser
	datetime          datetimes.Now
	maxUserAgentBytes int
}

func New(pool *pgxpool.Pool, client DBClient, idTool id.UUID7GeneratorParser, datetimeTool datetimes.Now, logger logs.Logger) (*Repo, error) {
	if pool == nil {
		return nil, errNilDBPool
	}

	if client == nil {
		return nil, errNilDBClient
	}

	if idTool == nil {
		return nil, errNilIDTool
	}

	if logger == nil {
		return nil, errNilLogger
	}

	if datetimeTool == nil {
		return nil, errNilDatetimeTool
	}

	return &Repo{
		pool:              pool,
		client:            client,
		id:                idTool,
		logger:            logger.With("component", "auth_repo"),
		datetime:          datetimeTool,
		maxUserAgentBytes: 512, // set a default value for maxUserAgentBytes
	}, nil
}
