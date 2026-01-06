package logs

import (
	"context"
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
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

func NewMultiSlogger(handlers ...slog.Handler) Logger {
	return &Slogger{slog.New(slogmulti.Fanout(handlers...))}
}

func (sl *Slogger) With(args ...any) Logger {
	return &Slogger{sl.Logger.With(args...)}
}

func (sl *Slogger) Log(ctx context.Context, level Level, msg string, args ...any) {
	sl.Logger.Log(ctx, slog.Level(level), msg, args...)
}

func (sl *Slogger) FatalContext(ctx context.Context, msg string, args ...any) {
	sl.Logger.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}
