package fighters

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/fighters"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type FighterGetter interface {
	GetFighterByID(ctx context.Context, id uuid.UUID) (service.Fighter, error)
}

// getFighter handles the HTTP POST request for the v1 create_Fighter endpoint.
func (h *Handler) getFighter(svc FighterGetter, logger logs.Logger) http.HandlerFunc {
	logger = logger.With(logs.FieldEndpoint, v1.EndpointNameV1UsersGet)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fighterIDStr := mux.Vars(r)[v1.QueryParamFighterID]
		fighterID, err := uuid.FromString(fighterIDStr)
		if err != nil {
			logger.DebugContext(ctx, "invalid fighter_id parameter", "error", err, "fighter_id", fighterIDStr)
			h.writeError(ctx, w, errors.Wrap(v1.ErrInvalidUUID, v1.QueryParamFighterID), logger)
			return
		}

		createdFighter, err := svc.GetFighterByID(ctx, fighterID)
		if err != nil {
			logger.ErrorContext(ctx, "failed to get fighter", "error", err)
			h.writeError(ctx, w, err, logger)
			return
		}

		resp := fighterIDOToDTO(createdFighter)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.ErrorContext(ctx, "failed to write response body", "error", err)
		}
	}
}
