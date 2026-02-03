package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type DBClient interface {
	GetUserRolesByID(ctx context.Context, userID id.UUID7) ([]string, error)
	GetUserPermissionsByID(ctx context.Context, id id.UUID7) ([]postgres.GetUserPermissionsByIDRow, error)
	StoreRefreshToken(ctx context.Context, arg postgres.StoreRefreshTokenParams) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (postgres.RefreshToken, error)
	RevokeRefreshTokenByHash(ctx context.Context, arg postgres.RevokeRefreshTokenByHashParams) error
	RevokeAndRotateRefreshTokenByHash(ctx context.Context, arg postgres.RevokeAndRotateRefreshTokenByHashParams) error
	WithTx(tx pgx.Tx) *postgres.Queries
}

type Repo struct {
	pool   *pgxpool.Pool
	client DBClient
	logger logs.Logger
	id     id.UUID7GeneratorParser
}

func New(pool *pgxpool.Pool, client DBClient, idTool id.UUID7GeneratorParser, logger logs.Logger) (*Repo, error) {
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

	return &Repo{
		pool:   pool,
		client: client,
		id:     idTool,
		logger: logger.With("component", "auth_repo"),
	}, nil
}
