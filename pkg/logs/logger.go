package logs

import (
	"context"
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

var (
	FieldEndpoint  = "endpoint"
	FieldRequestID = "request_id"
	FieldUserID    = "user_id"
)

type Level slog.Level

var (
	LevelDebug Level = Level(slog.LevelDebug)
	LevelInfo  Level = Level(slog.LevelInfo)
	LevelWarn  Level = Level(slog.LevelWarn)
	LevelError Level = Level(slog.LevelError)
	LevelFatal Level = Level(slog.LevelError) // slog does not have Fatal level
)

type Logger interface {
	DebugContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	FatalContext(ctx context.Context, msg string, args ...any)
	Log(ctx context.Context, level Level, msg string, args ...any)
	With(args ...any) Logger
}

type Slogger struct {
	*slog.Logger
}

// New creates a new Logger instance with the specified log level.
// If otelEnabled is true, an OTel slog handler is added for exporting logs to OpenTelemetry on top of the standard output JSON handler.
// Each handler is wrapped with ContextHandler to enrich logging calls with request-scoped context values.
func New(level Level, otelEnabled bool, appName, env string) Logger {
	loggerHandlers := []slog.Handler{NewContextHandler(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(level)}),
	)}

	if otelEnabled {
		loggerHandlers = append(loggerHandlers, NewContextHandler(otelslog.NewHandler(appName)))
	}

	// TODO: slogmulti can be removed once we have go 1.26 is released as it includes built-in support for multiple handlers https://tip.golang.org/doc/go1.26
	baseLogger := newMultiSlogger(loggerHandlers...)

	return baseLogger.With("env", env)
}

func newMultiSlogger(handlers ...slog.Handler) Logger {
	return &Slogger{slog.New(slogmulti.Fanout(handlers...))}
}

func (sl *Slogger) With(args ...any) Logger {
	return &Slogger{sl.Logger.With(args...)}
}

func (sl *Slogger) InfoContext(ctx context.Context, msg string, args ...any) {
	sl.Logger.InfoContext(ctx, msg, args...)
}

func (sl *Slogger) DebugContext(ctx context.Context, msg string, args ...any) {
	sl.Logger.DebugContext(ctx, msg, args...)
}

func (sl *Slogger) WarnContext(ctx context.Context, msg string, args ...any) {
	sl.Logger.WarnContext(ctx, msg, args...)
}

func (sl *Slogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	sl.Logger.ErrorContext(ctx, msg, args...)
}

func (sl *Slogger) Log(ctx context.Context, level Level, msg string, args ...any) {
	sl.Logger.Log(ctx, slog.Level(level), msg, args...)
}

func (sl *Slogger) FatalContext(ctx context.Context, msg string, args ...any) {
	sl.Logger.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}
