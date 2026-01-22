package fighters

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/OJOMB/fightpicker/internal/http/apiresponder"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

const pathPrefix = "/api/v1/fighters"

type Service interface {
	FighterGetter
	FightersIngestor
}

type Handler struct {
	*v1.Handler
	service    Service
	pathPrefix string
}

func New(service Service, idTool id.UUID7Parser, ctxTool contextual.ContextProvider, logger logs.Logger) (*Handler, error) {
	if logger == nil {
		return nil, v1.ErrLoggerIsNil
	}

	if idTool == nil {
		return nil, v1.ErrIDToolIsNil
	}

	if ctxTool == nil {
		return nil, v1.ErrContextToolIsNil
	}

	responder := apiresponder.NewJSONResponder(ctxTool, classifyError, logger.With("component", "handler_fighters_v1"))
	return &Handler{
		Handler:    v1.NewHandler(idTool, responder),
		service:    service,
		pathPrefix: pathPrefix,
	}, nil
}

func (h *Handler) RegisterRoutes(m *mux.Router) {
	// GET /api/v1/fighters/{fighter_id} - get a fighter by ID
	h.AddRoute(m, fmt.Sprintf("%s/{%s}", h.pathPrefix, v1.QueryParamFighterID), http.MethodGet, v1.EndpointNameV1FightersGet, h.getFighter(h.service))

	// POST /api/v1/fighters:ingest - ingest fighters
	h.AddRoute(m, fmt.Sprintf("%s:ingest", h.pathPrefix), http.MethodPost, v1.EndpointNameV1FightersIngest, h.ingestFighters(h.service))
}
