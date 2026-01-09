package fighters

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

func (s *Service) GetFighterByID(ctx context.Context, id uuid.UUID) (Fighter, error) {
	if id == uuid.Nil {
		return Fighter{}, ErrInvalidFighterID
	}

	fighter, err := s.repo.GetFighterByID(ctx, id)
	if err != nil {
		return Fighter{}, err
	}

	return fighter, nil
}
