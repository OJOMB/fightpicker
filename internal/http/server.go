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
	"github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type jwtValidator interface {
	Parse(tokenStr string, secretKey []byte) (*auth.AuthClaims, *jwt.RegisteredClaims, error)
}

// RouteRegistrar defines something that can register routes on a mux.
type RouteRegistrar interface {
	RegisterRoutes(mux *mux.Router, logger logs.Logger)
}

type Server struct {
	handlers     []RouteRegistrar
	router       *mux.Router
	addr         net.Addr
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

	handlers := []RouteRegistrar{
		handlersv1auth.New(app.Services.AuthService),
		handlersv1users.New(app.Services.UsersService),
		handlersv1fighters.New(app.Services.FightersService),
	}

	return &Server{
		handlers:     handlers,
		router:       router,
		addr:         &net.TCPAddr{IP: net.ParseIP(cfg.HTTP.Domain), Port: cfg.HTTP.Port},
		jwtValidator: app.Utils.JWTTool,
		logger:       app.Logger,
		oTelEnabled:  cfg.Observability.OTel.Enable,
		secretKey:    []byte(cfg.Auth.PrivateKey),
		env:          cfg.Env,
	}, nil
}

func (s *Server) WithHandlers(handlers []RouteRegistrar) *Server {
	s.handlers = handlers
	return s
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
