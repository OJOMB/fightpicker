package fighters

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/fighters"
)

type FighterGetter interface {
	GetFighterByID(ctx context.Context, id uuid.UUID) (service.Fighter, error)
}

// getFighter handles the HTTP POST request for the v1 create_Fighter endpoint.
func (h *Handler) getFighter(svc FighterGetter) v1.AppHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		fighterID, err := h.ParseID(r, v1.QueryParamFighterID)
		if err != nil {
			return err
		}

		createdFighter, err := svc.GetFighterByID(ctx, fighterID)
		if err != nil {
			return err
		}

		h.WriteJSON(ctx, w, http.StatusOK, fighterIDOToDTO(createdFighter))

		return nil
	}
}
