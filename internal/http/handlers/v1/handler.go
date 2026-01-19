package v1

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/OJOMB/fightpicker/internal/http/apierr"
	"github.com/OJOMB/fightpicker/internal/http/apiresponder"
	"github.com/OJOMB/fightpicker/pkg/id"
)

// HandlerFunc defines a function signature that returns a standard error.
// the error is intended to be classified and handled by the caller.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// APIErrClassifier is a function that classifies an error into an apierr.APIError.
type APIErrClassifier func(error) apierr.APIError

// Handler is the base HTTP handler for v1 endpoints.
type Handler struct {
	apiresponder.Responder
	id id.UUID7Parser
}

func NewHandler(idTool id.UUID7Parser, responder apiresponder.Responder) *Handler {
	return &Handler{
		Responder: responder,
		id:        idTool,
	}
}

func (h *Handler) AddRoute(m *mux.Router, OpPath, OpMethod, OpName string, hf HandlerFunc) {
	m.Handle(
		OpPath,
		otelhttp.NewHandler(h.ToHandler(hf), OpName),
	).Name(OpName).
		Methods(OpMethod)
}

func (h *Handler) ToHandler(action HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := action(w, r); err != nil {
			h.WriteError(r.Context(), w, err)
		}
	}
}

// ParsePaginationParams is a helper that parses pagination parameters from the HTTP request.
func (h *Handler) ParsePaginationParams(r *http.Request) (pageSize int, lastSeenID *id.UUID7, err error) {
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
func (h *Handler) parseLastSeenID(r *http.Request) (*id.UUID7, error) {
	lastSeenIDStr := r.URL.Query().Get(QueryParamLastSeenID)
	if lastSeenIDStr == "" {
		return nil, nil
	}

	lastSeenID, err := h.id.ParseString(lastSeenIDStr)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidUUID, QueryParamLastSeenID)
	}

	return &lastSeenID, nil
}

// ParseID is a helper that parses a UUID from the request URL path variables.
func (h *Handler) ParseID(r *http.Request, queryParam string) (id.UUID7, error) {
	idStr := mux.Vars(r)[queryParam]
	userID, err := h.id.ParseString(idStr)
	if err != nil {
		return id.UUID7Nil, errors.Wrap(ErrInvalidUUID, queryParam)
	}

	return userID, nil
}
