package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/OJOMB/fightpicker/internal/app"
	"github.com/OJOMB/fightpicker/internal/config"
	userspb "github.com/OJOMB/fightpicker/internal/grpc/gen/go/users/v1"
	"github.com/OJOMB/fightpicker/internal/grpc/interceptors"
	usersgrpc "github.com/OJOMB/fightpicker/internal/grpc/users"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type Server struct {
	grpcServer *grpc.Server
	lis        net.Listener
	Logger     logs.Logger
}

func New(cfg *config.Config, app *app.App) (*Server, error) {
	if cfg == nil {
		return nil, errNilConfig
	}

	if app == nil {
		return nil, errNilApp
	}

	interceptAuth, err := interceptors.NewUnaryAuthInterceptor([]byte(cfg.Auth.PrivateKey), app.Utils.JWTTool, app.Utils.ContextTool)
	if err != nil {
		return nil, errors.Wrap(errFailedToInitInterceptor, err.Error())
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(interceptAuth.GetInterceptor()),
	)

	usersServer := usersgrpc.NewServer(
		app.Services.UsersService,
		&app.Logger,
	)

	userspb.RegisterUsersServiceServer(grpcServer, usersServer)

	if cfg.GRPC.EnableReflection {
		reflection.Register(grpcServer)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GRPC.Domain, cfg.GRPC.Port))
	if err != nil {
		return nil, errors.Wrapf(errListeningFailed, "'%s:%d': %s", cfg.GRPC.Domain, cfg.GRPC.Port, err.Error())
	}

	return &Server{
		grpcServer: grpcServer,
		lis:        lis,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	s.Logger.InfoContext(ctx, "starting gRPC server", "addr", s.lis.Addr().String())
	return s.grpcServer.Serve(s.lis)
}
