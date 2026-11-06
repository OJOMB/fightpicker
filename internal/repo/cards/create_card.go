package cards

import (
	"context"

	service "github.com/OJOMB/fightpicker/internal/service/cards"
)

// CreateCard converts a Card IDO to a DBO and calls the repo to create the card in the database
func (r *Repo) CreateCard(ctx context.Context, user service.Card) error {
	if err := r.dbClient.CreateCard(ctx, cardIDOtoDBO(user)); err != nil {
		return dbErrorToServiceError(err)
	}

	return nil
}
