package middlewares

import (
	"net/http"

	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type ContextLoader struct {
	ctxTool contextual.ContextRequestProvider
	logger  logs.Logger
}

func NewContextLoader(ctxTool contextual.ContextRequestProvider, logger logs.Logger) *ContextLoader {
	return &ContextLoader{
		ctxTool: ctxTool,
		logger:  logger.With("middleware", "ContextLoader"),
	}
}

func (cl *ContextLoader) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := cl.ctxTool.WithRequestValues(r.Context(), r)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
