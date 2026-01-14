package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/OJOMB/fightpicker/internal/http/apierr"
	"github.com/OJOMB/fightpicker/internal/http/dtos"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// HandlerFunc defines a function signature that returns a standard error.
// the error is intended to be classified and handled by the caller.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// APIErrClassifier is a function that classifies an error into an apierr.APIError.
type APIErrClassifier func(error) apierr.APIError

// Handler is the base HTTP handler for v1 endpoints.
type Handler struct {
	Logger logs.Logger
}

func (h *Handler) WriteJSON(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	v any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.Logger.ErrorContext(ctx, "failed to write response body", "error", err)
	}
}

func (h *Handler) ToHandler(action HandlerFunc, classifier APIErrClassifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := action(w, r); err != nil {
			h.WriteError(r.Context(), w, classifier(err))
		}
	}
}

// writeError is a helper function to create a JSON formatted error from an api level error.
// It maps specific errors to appropriate HTTP status codes and logs the errors accordingly.
// log level is determined based on the severity of the error. Error logs are reserved for genuine server-side issues,
// while client-side errors are logged at the debug level.
func (h *Handler) WriteError(ctx context.Context, w http.ResponseWriter, err apierr.APIError) {
	reqID, _ := ctx.Value(contextual.KeyRequestID).(string)
	if reqID == "" {
		h.Logger.WarnContext(ctx, "request ID missing from context")
		reqID = "unknown"
	}

	switch err.LogLevel {
	case logs.LevelDebug:
		h.Logger.DebugContext(ctx, err.LogMsg, "error", err)
	case logs.LevelError:
		h.Logger.ErrorContext(ctx, err.LogMsg, "error", err)
	}

	resp := dtos.NewErrorEnvelope(err.Public, err.Code, reqID)
	h.WriteJSON(ctx, w, err.Status, resp)
}

// ParsePaginationParams is a helper that parses pagination parameters from the HTTP request.
func (h *Handler) ParsePaginationParams(r *http.Request) (pageSize int, lastSeenID *uuid.UUID, err error) {
	pageSize, err = h.parsePageSize(r)
	if err != nil {
		return 0, nil, err
	}

	lastSeenID, err = h.parseLastSeenID(r)
	if err != nil {
		return 0, nil, err
	}

	return pageSize, lastSeenID, nil
}

// parsePageSize is a pagination helper that parses the page size from the query parameter string.
func (h *Handler) parsePageSize(r *http.Request) (int, error) {
	pageSizeStr := r.URL.Query().Get(QueryParamPageSize)

	if pageSizeStr == "" {
		return DefaultPageSize, nil
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 0 {
		return 0, errors.Wrap(ErrInvalidPageSize, QueryParamPageSize)
	}

	if pageSize > MaxPageSize {
		return MaxPageSize, nil
	}

	return pageSize, nil
}

// parseLastSeenID is a pagination helper that parses the last seen ID from the query parameter string.
func (h *Handler) parseLastSeenID(r *http.Request) (*uuid.UUID, error) {
	lastSeenIDStr := r.URL.Query().Get(QueryParamLastSeenID)
	if lastSeenIDStr == "" {
		return nil, nil
	}

	id, err := uuid.FromString(lastSeenIDStr)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidUUID, QueryParamLastSeenID)
	}

	return &id, nil
}

// ParseID is a helper that parses a UUID from the request URL path variables.
func (h *Handler) ParseID(r *http.Request, queryParam string) (uuid.UUID, error) {
	idStr := mux.Vars(r)[queryParam]
	userID, err := uuid.FromString(idStr)
	if err != nil {
		return uuid.Nil, errors.Wrap(ErrInvalidUUID, queryParam)
	}

	return userID, nil
}
