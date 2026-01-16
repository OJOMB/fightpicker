package fighters

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

const pathPrefix = "/api/v1/fighters"

type Service interface {
	FighterGetter
}

type Handler struct {
	*v1.Handler
	service    Service
	pathPrefix string
}

func New(service Service, idTool id.UUID7Parser, logger logs.Logger) *Handler {
	return &Handler{
		Handler:    v1.NewHandler(idTool, logger),
		service:    service,
		pathPrefix: pathPrefix,
	}
}

func (h *Handler) RegisterRoutes(mux *mux.Router) {
	// GET /api/v1/fighters/{fighter_id} - get a fighter by ID
	mux.Handle(
		fmt.Sprintf("%s/{%s}", h.pathPrefix, v1.QueryParamFighterID),
		otelhttp.NewHandler(
			h.ToHandler(
				h.getFighter(h.service),
				classifyError,
			),
			v1.EndpointNameV1FightersGet,
		),
	).Name(v1.EndpointNameV1FightersGet).
		Methods(http.MethodGet)
}
