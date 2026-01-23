package fighters

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	fightersservice "github.com/OJOMB/fightpicker/internal/service/fighters"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

func (r *Repo) IngestFighters(ctx context.Context, rows []fightersservice.IngestRow, adminUserID id.UUID7) ([]fightersservice.IngestionResult, fightersservice.IngestionSummary, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fightersservice.IngestionSummary{}, err
	}
	defer tx.Rollback(ctx)

	qs := r.client.WithTx(tx)

	// Convert rows → JSON payload expected by SQL
	jsonPayload, err := json.Marshal(rows)
	if err != nil {
		return nil, fightersservice.IngestionSummary{}, err
	}

	params := postgres.IngestFightersParams{
		Payload: jsonPayload,
		OperationTime: pgtype.Timestamptz{
			Time:  r.now.Now().UTC(),
			Valid: true,
		},
		AdminUserID: pgtype.UUID{
			Bytes: adminUserID,
			Valid: true,
		},
	}

	dbRows, err := qs.IngestFighters(ctx, params)
	if err != nil {
		return nil, fightersservice.IngestionSummary{}, dbErrorToServiceError(err)
	}

	results := make([]fightersservice.IngestionResult, 0, len(dbRows))
	var summary fightersservice.IngestionSummary

	for _, row := range dbRows {
		result := fightersservice.IngestionResult{
			Index:  int(row.Idx),
			Status: fightersservice.IngestionResultStatus(row.Status),
		}

		if row.FighterID != id.UUID7Nil {
			result.FighterId = &row.FighterID
		}

		if row.ErrorCode.Valid {
			result.Error = &fightersservice.ErrorObject{
				Code:    row.ErrorCode.String,
				Message: row.ErrorMessage.String,
			}
		}

		results = append(results, result)

		// Update summary
		summary.Total++
		switch fightersservice.IngestionResultStatus(row.Status) {
		case fightersservice.IngestionResultStatusCreated:
			summary.Created++
		case fightersservice.IngestionResultStatusUpdated:
			summary.Updated++
		case fightersservice.IngestionResultStatusSkipped:
			summary.Skipped++
		case fightersservice.IngestionResultStatusFailed:
			summary.Failed++
		default:
			r.logger.ErrorContext(ctx, "unknown ingestion result status", "status", row.Status)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fightersservice.IngestionSummary{}, dbErrorToServiceError(err)
	}

	return results, summary, nil
}
