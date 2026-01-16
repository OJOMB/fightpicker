package fighters

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/id"
)

func (s *Service) GetFighterByID(ctx context.Context, fighterID id.UUID7) (Fighter, error) {
	if fighterID == id.UUID7Nil {
		return Fighter{}, ErrInvalidFighterID
	}

	fighter, err := s.repo.GetFighterByID(ctx, fighterID)
	if err != nil {
		return Fighter{}, err
	}

	return fighter, nil
}
