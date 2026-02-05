package fighters

import (
	"context"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/fighters"
	"github.com/OJOMB/fightpicker/pkg/id"
)

type FighterLister interface {
	ListFighters(ctx context.Context, pageSize int, lastSeenID *id.UUID7) ([]service.Fighter, int, error)
}

type FighterSearcher interface {
	FighterLister
}

// listFighters handles the HTTP GET request for the v1 list_users endpoint.
func (h *Handler) listFighters(svc FighterSearcher) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		pageSize, lastSeenID, err := h.ParsePaginationParams(r)
		if err != nil {
			return err
		}

		fighters, totalCount, err := svc.ListFighters(ctx, pageSize, lastSeenID)
		if err != nil {
			return err
		}

		resp := dtos.ListFightersResponse{
			TotalCount: totalCount,
			Fighters:   make([]dtos.FighterResponse, len(fighters)),
			PageSize:   len(fighters),
		}
		for i, fighter := range fighters {
			resp.Fighters[i] = fighterIDOToDTO(fighter)
		}

		// set LastSeenId for pagination if there are more users to fetch
		if len(fighters) > 0 && len(fighters) == pageSize {
			resp.LastSeenId = &resp.Fighters[len(resp.Fighters)-1].Id
		}

		h.Write(ctx, w, http.StatusOK, resp)
		return nil
	}
}
