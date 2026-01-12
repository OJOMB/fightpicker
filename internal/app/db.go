package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/multitracer"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"

	"github.com/OJOMB/fightpicker/internal/config"
	postgresclient "github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type db struct {
	pool    *pgxpool.Pool
	queries *postgresclient.Queries
}

func (a *App) newDB(ctx context.Context, cfg *config.Config) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Auth.SSLMode,
	)

	// 1. Create pool config
	dbCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}

	// create multitracer for otelsql
	tracer := multitracer.New(otelpgx.NewTracer(), &tracelog.TraceLog{
		Logger:   SlogTracer{a.Logger},
		LogLevel: tracelog.LogLevelInfo,
	})

	dbCfg.ConnConfig.Tracer = tracer

	// 3. Create the pool (thread-safe)
	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		return err
	}

	// 4. Inject pool into sqlc generated client
	q := postgresclient.New(pool)

	a.DB = db{
		pool:    pool,
		queries: q,
	}

	return nil
}

// SlogTracer implements pgx tracelog.Logger using our slog.Logger
type SlogTracer struct {
	logger logs.Logger
}

// Log implements pgx tracelog.Logger converting it to a regular slog log
func (s SlogTracer) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	// Convert map[string]any → slog style attributes
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	s.logger.Log(ctx, logs.Level(level), msg, attrs)
}
