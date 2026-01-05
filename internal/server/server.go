package server

import (
	"context"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"

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

func New(domain string, port int, router *mux.Router, jwtValidator jwtValidator, secretKey string, oTelEnabled bool, env string, logger logs.Logger) (*Server, error) {
	if router == nil {
		router = mux.NewRouter()
	}

	if domain == "" {
		return nil, ErrNoDomainSpecified
	}

	if port == 0 {
		return nil, ErrPortNotSpecified
	}

	return &Server{
		router:       router,
		addr:         &net.TCPAddr{IP: net.ParseIP(domain), Port: port},
		jwtValidator: jwtValidator,
		logger:       logger,
		oTelEnabled:  oTelEnabled,
		secretKey:    []byte(secretKey),
		env:          env,
	}, nil
}

func (svr *Server) WithHandlers(handlers []RouteRegistrar) *Server {
	svr.handlers = handlers
	return svr
}

// Run starts the HTTP server and listens for incoming requests. It also handles graceful shutdown on receiving termination signals.
func (svr *Server) Run() error {
	svr.routes()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	httpServer := &http.Server{
		Handler: svr.router,
	}

	// Create listener
	listener, err := net.Listen("tcp", svr.addr.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	svr.logger.InfoContext(ctx, "starting server", "domain", svr.addr.String())

	// run http server in background
	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.Serve(listener)
	}()

	// wait for signal or server error
	select {
	case <-ctx.Done():
		svr.logger.InfoContext(ctx, "shutdown signal received")
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
