package main

import (
	"context"
	"flag"
	"log"

	"github.com/OJOMB/fightpicker/internal/app"
	"github.com/OJOMB/fightpicker/internal/config"
	"github.com/OJOMB/fightpicker/internal/grpc"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/pyroscope"
)

var env string

func main() {
	flag.StringVar(&env, "env", "local", "runtime environment (local|e2e|staging|prod)")
	flag.Parse()

	cfg, err := config.Load(env)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	ctx, stop := contextual.SetupSignals()
	defer stop()

	if cfg.Observability.Pyroscope.Enable {
		pyroProf, err := pyroscope.Setup(cfg.AppName, cfg.Observability.Pyroscope.Endpoint)
		if err != nil {
			log.Fatalf("failed to start pyroscope: %v", err)
		}

		defer pyroProf.Stop()
	}

	app, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer app.Shutdown(context.Background())

	srv, err := grpc.New(cfg, app)
	if err != nil {
		app.Logger.FatalContext(ctx, "failed to create gRPC server", "error", err)
	}

	if err := srv.Run(ctx); err != nil {
		app.Logger.FatalContext(ctx, "gRPC server encountered an error", "error", err)
	}
}
