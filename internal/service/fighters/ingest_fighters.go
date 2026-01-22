package fighters

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/id"
)

type FightersIngestor interface {
	IngestFighters(ctx context.Context, ingestRow []IngestRow, adminUUID id.UUID7) ([]IngestionResult, IngestionSummary, error)
}

func (s *Service) IngestFighters(ctx context.Context, fighterIngestionRows []IngestRow) ([]IngestionResult, IngestionSummary, error) {
	reqSubject := s.ctxTool.GetReqSubjectFromContext(ctx)
	if reqSubject == id.UUID7Nil || !s.ctxTool.ReqSubjectIsAnAdmin(ctx) {
		return nil, IngestionSummary{}, ErrUnauthorized
	}

	if len(fighterIngestionRows) == 0 {
		return nil, IngestionSummary{}, nil
	}

	results := make([]IngestionResult, len(fighterIngestionRows))
	valid := make([]IngestRow, 0, len(fighterIngestionRows))

	var failed int

	// Pre-fill result indices
	for i, f := range fighterIngestionRows {
		results[i].Index = f.Index
	}

	// Validate
	for i, f := range fighterIngestionRows {
		if err := s.validateCreationReq(&f.Fighter); err != nil {
			results[i].Status = IngestionResultStatusFailed
			results[i].Error = &ErrorObject{
				Code:    "validation_error",
				Message: err.Error(),
			}
			failed++

			continue
		}

		valid = append(valid, f)
	}

	// Nothing valid to ingest
	if len(valid) == 0 {
		return results, IngestionSummary{
			Failed: failed,
			Total:  len(fighterIngestionRows),
		}, nil
	}

	// DB ingest
	ingestResults, ingestSummary, err := s.repo.IngestFighters(ctx, valid, reqSubject)
	if err != nil {
		return nil, IngestionSummary{}, err
	}

	// Merge DB results back
	dbIdx := 0
	for i := range results {
		if results[i].Status == IngestionResultStatusFailed {
			continue
		}

		results[i] = ingestResults[dbIdx]
		dbIdx++
	}

	// Final summary
	ingestSummary.Failed += failed
	ingestSummary.Total = len(fighterIngestionRows)

	return results, ingestSummary, nil
}
