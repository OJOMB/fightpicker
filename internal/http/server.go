package http

import (
	"context"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"

	"github.com/OJOMB/fightpicker/internal/app"
	"github.com/OJOMB/fightpicker/internal/config"
	handlersv1auth "github.com/OJOMB/fightpicker/internal/http/handlers/v1/auth"
	handlersv1fighters "github.com/OJOMB/fightpicker/internal/http/handlers/v1/fighters"
	handlersv1users "github.com/OJOMB/fightpicker/internal/http/handlers/v1/users"
	"github.com/OJOMB/fightpicker/internal/http/middlewares"
	"github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type jwtValidator interface {
	Parse(tokenStr string, secretKey []byte) (*auth.AuthClaims, *jwt.RegisteredClaims, error)
}

// RouteRegistrar defines something that can register routes on a mux.
type RouteRegistrar interface {
	RegisterRoutes(mux *mux.Router)
}

type Server struct {
	handlers     []RouteRegistrar
	middlewares  []mux.MiddlewareFunc
	router       *mux.Router
	addr         net.Addr
	id           id.UUID7GeneratorParser
	jwtValidator jwtValidator
	secretKey    []byte
	logger       logs.Logger
	oTelEnabled  bool
	env          string
}

func New(cfg *config.Config, app *app.App) (*Server, error) {
	router := mux.NewRouter()

	if cfg.HTTP.Domain == "" {
		return nil, ErrNoDomainSpecified
	}

	if cfg.HTTP.Port == 0 {
		return nil, ErrPortNotSpecified
	}

	handlerAuth, err := handlersv1auth.New(app.Services.AuthService, app.Utils.IDTool, app.Utils.ContextTool, app.Logger)
	if err != nil {
		return nil, err
	}

	handlerUsers, err := handlersv1users.New(app.Services.UsersService, app.Utils.IDTool, app.Utils.ContextTool, app.Logger)
	if err != nil {
		return nil, err
	}

	handlerFighters, err := handlersv1fighters.New(app.Services.FightersService, app.Utils.IDTool, app.Utils.ContextTool, app.Logger)
	if err != nil {
		return nil, err
	}

	handlers := []RouteRegistrar{
		handlerAuth,
		handlerUsers,
		handlerFighters,
	}

	authMW, err := middlewares.NewAuthPermissionsChecker([]byte(cfg.Auth.PrivateKey), app.Utils.JWTTool, app.Utils.ContextTool, app.Logger)

	logRespBody := cfg.Env != "production"
	middlewares := []mux.MiddlewareFunc{
		middlewares.NewRequestResponseLogger(app.Logger, app.Utils.IDTool, logRespBody, cfg.Observability.OTel.Enable).Middleware,
		authMW.Middleware,
		middlewares.NewContextLoader(app.Utils.ContextTool, app.Logger).Middleware,
		middlewares.NewPyroProfiler(map[string]string{"component": "server"}).Middleware,
	}

	return &Server{
		handlers:     handlers,
		middlewares:  middlewares,
		router:       router,
		addr:         &net.TCPAddr{IP: net.ParseIP(cfg.HTTP.Domain), Port: cfg.HTTP.Port},
		jwtValidator: app.Utils.JWTTool,
		id:           app.Utils.IDTool,
		logger:       app.Logger,
		oTelEnabled:  cfg.Observability.OTel.Enable,
		secretKey:    []byte(cfg.Auth.PrivateKey),
		env:          cfg.Env,
	}, nil
}

// Run starts the HTTP server and listens for incoming requests. It also handles graceful shutdown on receiving termination signals.
func (s *Server) Run() error {
	s.routes()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	httpServer := &http.Server{
		Handler: s.router,
	}

	// Create listener
	listener, err := net.Listen("tcp", s.addr.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	s.logger.InfoContext(ctx, "starting server", "domain", s.addr.String())

	// run http server in background
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.Serve(listener)
	}()

	// wait for signal or server error
	select {
	case <-ctx.Done():
		s.logger.InfoContext(ctx, "shutdown signal received")
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	// commence graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
