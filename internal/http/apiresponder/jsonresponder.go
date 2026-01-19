package apiresponder

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/apierr"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

const (
	headerContentTypeKey       = "Content-Type"
	contentTypeApplicationJSON = "application/json"
)

type JSONResponder struct {
	ctxTool         contextual.ContextProvider
	errorClassifier apierr.APIErrClassifier
	logger          logs.Logger
}

func NewJSONResponder(
	ctxTool contextual.ContextProvider,
	errorClassifier apierr.APIErrClassifier,
	logger logs.Logger,
) *JSONResponder {
	return &JSONResponder{
		ctxTool:         ctxTool,
		errorClassifier: errorClassifier,
		logger:          logger,
	}
}

func (r *JSONResponder) WriteError(
	ctx context.Context,
	w http.ResponseWriter,
	err error,
) {
	reqID := r.ctxTool.GetRequestID(ctx)
	if reqID == id.UUID7Nil {
		r.logger.WarnContext(ctx, "request ID missing from context")
	}

	apiErr := r.errorClassifier(err)

	switch apiErr.LogLevel {
	case logs.LevelDebug:
		r.logger.DebugContext(ctx, apiErr.LogMsg, "error", err)
	case logs.LevelError:
		r.logger.ErrorContext(ctx, apiErr.LogMsg, "error", err)
	}

	r.Write(ctx, w, apiErr.Status, apiErr.ToDTO(reqID))
}

func (r *JSONResponder) Write(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	v any,
) {
	w.Header().Set(headerContentTypeKey, contentTypeApplicationJSON)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		r.logger.ErrorContext(ctx, "failed to write response body", "error", err)
	}
}
