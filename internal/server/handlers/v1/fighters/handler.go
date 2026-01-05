package fighters

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

const pathPrefix = "/api/v1/fighters"

type Service interface {
	FighterGetter
}

type Handler struct {
	service    Service
	pathPrefix string
}

func New(service Service) *Handler {
	return &Handler{
		service:    service,
		pathPrefix: pathPrefix,
	}
}

func (h *Handler) RegisterRoutes(mux *mux.Router, logger logs.Logger) {
	// // POST /api/v1/users - create a new user
	// mux.Handle(
	// 	h.pathPrefix,
	// 	otelhttp.NewHandler(
	// 		h.createUser(h.service, logger),
	// 		v1.EndpointNameV1UsersCreate,
	// 	),
	// ).Name(v1.EndpointNameV1UsersCreate).
	// 	Methods(http.MethodPost)

	// GET /api/v1/fighters/{fighter_id} - get a fighter by ID
	mux.Handle(
		fmt.Sprintf("%s/{%s}", h.pathPrefix, v1.QueryParamFighterID),
		otelhttp.NewHandler(
			h.getFighter(h.service, logger),
			v1.EndpointNameV1FightersGet,
		),
	).Name(v1.EndpointNameV1FightersGet).
		Methods(http.MethodGet)
}
