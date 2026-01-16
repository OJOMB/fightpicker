package logs

import (
	"context"
	"log/slog"
	"net/http"
)

type ContextRequestProvider interface {
	WithRequestValues(ctx context.Context, r *http.Request) context.Context
	GetRequestValues(ctx context.Context) map[string]string
}

type ContextHandler struct {
	ctxTool ContextRequestProvider
	next    slog.Handler
}

func NewContextHandler(ctxTool ContextRequestProvider, next slog.Handler) slog.Handler {
	return &ContextHandler{ctxTool: ctxTool, next: next}
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// Pull values from context
	reqValues := h.ctxTool.GetRequestValues(ctx)
	attrs := make([]slog.Attr, 0, len(reqValues))
	for k, v := range reqValues {
		attrs = append(attrs, slog.String(k, v))
	}

	// Add attributes to the record
	r.AddAttrs(attrs...)

	return h.next.Handle(ctx, r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{ctxTool: h.ctxTool, next: h.next.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{ctxTool: h.ctxTool, next: h.next.WithGroup(name)}
}
