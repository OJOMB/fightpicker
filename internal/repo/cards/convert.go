package cards

import (
	"github.com/jackc/pgx/v5/pgtype"

	service "github.com/OJOMB/fightpicker/internal/service/cards"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

func cardIDOtoDBO(c service.Card) postgres.CreateCardParams {
	return postgres.CreateCardParams{
		ID:       c.ID,
		Name:     c.Name,
		Location: c.Location,
		Date: pgtype.Timestamptz{
			Time:  c.Date,
			Valid: !c.Date.IsZero(),
		},
		Promotion: c.Promotion,

		CreatedAt: pgtype.Timestamptz{
			Time:  c.CreatedAt,
			Valid: !c.CreatedAt.IsZero(),
		},
		CreatedBy: pgtype.UUID{
			Bytes: c.CreatedBy,
			Valid: c.CreatedBy != id.UUID7Nil,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  c.UpdatedAt,
			Valid: !c.UpdatedAt.IsZero(),
		},
		UpdatedBy: pgtype.UUID{
			Bytes: c.UpdatedBy,
			Valid: c.UpdatedBy != id.UUID7Nil,
		},
	}
}
