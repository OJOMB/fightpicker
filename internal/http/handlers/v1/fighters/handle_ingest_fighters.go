package fighters

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	fightersservice "github.com/OJOMB/fightpicker/internal/service/fighters"
)

type FightersIngestor interface {
	IngestFighters(ctx context.Context, fighters []fightersservice.IngestRow) ([]fightersservice.IngestionResult, fightersservice.IngestionSummary, error)
}

// ingestFighters handles the HTTP POST request for the v1 ingest endpoint.
func (h *Handler) ingestFighters(svc FightersIngestor) v1.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		defer r.Body.Close()
		var fightersBatchCreateReq dtos.FightersIngestionReq
		if err := json.NewDecoder(r.Body).Decode(&fightersBatchCreateReq); err != nil {
			return v1.ErrInvalidJSONRequestBody
		}

		results, summary, err := svc.IngestFighters(ctx, fightersIngestionReqDTOtoIDO(fightersBatchCreateReq))
		if err != nil {
			return err
		}

		resp := ingestionIDOstoDTO(results, summary)

		h.Write(ctx, w, http.StatusOK, resp)

		return nil
	}
}
