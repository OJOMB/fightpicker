package cards

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/OJOMB/fightpicker/internal/http/apiresponder"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

const pathPrefix = "/api/v1/cards"

type Service interface {
	CardCreator
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

	responder := apiresponder.NewJSONResponder(ctxTool, classifyError, logger.With("component", "handler_cards_v1"))
	return &Handler{
		Handler:    v1.NewHandler(idTool, responder),
		service:    service,
		pathPrefix: pathPrefix,
	}, nil
}

func (h *Handler) RegisterRoutes(m *mux.Router) {
	// POST /api/v1/cards - create card
	h.AddRoute(m, h.pathPrefix, http.MethodPost, v1.EndpointNameV1CardsCreate, h.createCard(h.service))
}
