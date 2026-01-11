package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"

	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

const reqIDHeader = "X-Request-ID"

type responseWriterRecorder struct {
	http.ResponseWriter
	status int
	body   string
}

func (w *responseWriterRecorder) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterRecorder) Write(bs []byte) (int, error) {
	// Default to 200 if WriteHeader hasn't been called
	if w.status == 0 {
		w.status = http.StatusOK
	}

	w.body += string(bs)
	return w.ResponseWriter.Write(bs)
}

type RequestResponseLogger struct {
	logger          logs.Logger
	logResponseBody bool
	oTelEnabled     bool
}

func NewRequestResponseLogger(l logs.Logger, logRespBody, oTelEnabled bool) *RequestResponseLogger {
	return &RequestResponseLogger{
		logger:          l.With("middleware", "RequestResponseLogger"),
		logResponseBody: logRespBody,
		oTelEnabled:     oTelEnabled,
	}
}

func (rrl *RequestResponseLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.Must(uuid.NewV4()).String()
			r.Header.Set(reqIDHeader, reqID)
		}

		ctx := r.Context()

		// PLACE REQID IN CONTEXT FOR LOGGING
		ctx = context.WithValue(ctx, contextual.KeyRequestID, reqID)
		r = r.WithContext(ctx)

		// assumuing here that the OTel middleware has already run and added the span to the context
		var traceID, spanID string
		if rrl.oTelEnabled {
			span := trace.SpanFromContext(r.Context())
			sc := span.SpanContext()

			traceID = sc.TraceID().String()
			spanID = sc.SpanID().String()
		}

		rrl.logger.InfoContext(ctx, "incoming request",
			slog.String("type", "request"),
			slog.String("path", r.URL.Path),
			slog.String("query_parameters", r.URL.RawQuery),
			slog.String("method", r.Method),
			slog.String("request_id", reqID),
			slog.String("src_host", r.RemoteAddr),
			slog.String("trace_id", traceID),
			slog.String("span_id", spanID),
		)

		start := time.Now().UTC()

		// add X-Request-ID to response headers
		w.Header().Set(reqIDHeader, reqID)

		customResponseWriter := &responseWriterRecorder{ResponseWriter: w}
		next.ServeHTTP(customResponseWriter, r)

		end := time.Now().UTC()
		rrl.logger.InfoContext(ctx, "outgoing response",
			slog.String("type", "response"),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("request_id", reqID),
			slog.Int("status_code", customResponseWriter.status),
			slog.String("src_host", r.RemoteAddr),
			slog.Duration("duration", end.Sub(start)),
			slog.String("response_body", func() string {
				if rrl.logResponseBody {
					return customResponseWriter.body
				}
				return "omitted"
			}()),
			slog.String("trace_id", traceID),
			slog.String("span_id", spanID),
		)
	})
}
