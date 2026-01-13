// app builds the World.
// It holds references to all services, repositories, clients, utilities, and database connections required by transport layers.
// It is designed to be instantiated in main and passed to whichever transport layers to handle incoming requests.
// Thus avoiding duplicated setup code across transport layers.
package app

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"

	"github.com/OJOMB/fightpicker/internal/config"
	"github.com/OJOMB/fightpicker/pkg/logs"
	"github.com/OJOMB/fightpicker/pkg/otel"
)

// App is the main application struct that holds references to services, repositories, clients, utilities, and database connections.
// App is responsible for consistent initialization and configuration of these components, regardless of the transport layer used (e.g., HTTP, gRPC).
type App struct {
	DB       db
	Clients  clients
	Repos    repos
	Services services
	Utils    *utils

	otelShutdown func(ctx context.Context) error
	initialized  bool

	Logger logs.Logger
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	var a = new(App)

	loggerHandlers := []slog.Handler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(cfg.LogLevel)})}
	if cfg.Observability.OTel.Enable {
		log.Print("OTel enabled")
		otelShutdown, err := otel.SetupOTelSDK(ctx, cfg.Observability.OTel.Endpoint, cfg.AppName)
		if err != nil {
			return nil, err
		}

		a.otelShutdown = otelShutdown

		loggerHandlers = append(loggerHandlers, otelslog.NewHandler(cfg.AppName))
	}

	// TODO: slogmulti can be removed once we have go 1.26 is released as it includes built-in support for multiple handlers https://tip.golang.org/doc/go1.26
	baseLogger := logs.NewMultiSlogger(loggerHandlers...)

	a.Logger = baseLogger.With("env", cfg.Env)

	// 1. Initialize DB
	if err := a.newDB(ctx, cfg); err != nil {
		return nil, err
	}

	// 2. Initialize common utilities
	if err := a.newUtils(cfg); err != nil {
		return nil, err
	}

	// 3. Initialize Clients
	if err := a.newClients(ctx, cfg); err != nil {
		return nil, err
	}

	// 4. Initialize Repositories
	if err := a.newRepos(ctx, cfg); err != nil {
		return nil, err
	}

	// 5. Initialize Services
	if err := a.newServices(ctx, cfg); err != nil {
		return nil, err
	}

	a.initialized = true

	return a, nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var errs []error

	if a.DB.pool != nil {
		a.DB.pool.Close()
	}

	if a.otelShutdown != nil {
		errs = append(errs, a.otelShutdown(ctx))
	}

	return errors.Join(errs...)
}
