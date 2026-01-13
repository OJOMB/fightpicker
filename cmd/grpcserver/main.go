package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/OJOMB/fightpicker/internal/app"
	"github.com/OJOMB/fightpicker/internal/config"
	userspb "github.com/OJOMB/fightpicker/internal/grpc/gen/go/users/v1"
	usersgrpc "github.com/OJOMB/fightpicker/internal/grpc/users"
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

	grpcServer := grpc.NewServer(
	// interceptors go here later
	)

	usersServer := usersgrpc.NewServer(
		app.Services.UsersService,
		&app.Logger,
	)

	userspb.RegisterUsersServiceServer(grpcServer, usersServer)

	if cfg.GRPC.EnableReflection {
		reflection.Register(grpcServer)
	}

	// 6. Listen
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GRPC.Domain, cfg.GRPC.Port))
	if err != nil {
		panic(err)
	}

	app.Logger.InfoContext(ctx, "gRPC server listening", "addr", fmt.Sprintf("%s:%d", cfg.GRPC.Domain, cfg.GRPC.Port))

	// 7. Serve (blocking)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
