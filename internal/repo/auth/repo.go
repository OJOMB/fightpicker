package auth

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type DBClient interface {
	GetUserRolesByID(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetUserPermissionsByID(ctx context.Context, id uuid.UUID) ([]postgres.GetUserPermissionsByIDRow, error)
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
	idGen  id.Generator
}

func New(pool *pgxpool.Pool, client DBClient, idGen id.Generator, logger logs.Logger) *Repo {
	return &Repo{
		pool:   pool,
		client: client,
		idGen:  idGen,
		logger: logger.With("component", "auth_repo"),
	}
}
