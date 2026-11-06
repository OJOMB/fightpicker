package cards

import (
	"context"
	"fmt"
)

// CreateCard validates and stores a new card
func (s *Service) CreateCard(ctx context.Context, c Card) (Card, error) {
	if c.Name == "" {
		return Card{}, fmt.Errorf("%w: name", ErrMissingParameter)
	}

	if c.Date.IsZero() {
		return Card{}, fmt.Errorf("%w: date", ErrMissingParameter)
	}

	if c.Location == "" {
		return Card{}, fmt.Errorf("%w: location", ErrMissingParameter)
	}

	if c.Promotion == "" {
		return Card{}, fmt.Errorf("%w: promotion", ErrMissingParameter)
	}

	// fill in missing fields
	c.ID = s.idTool.Generate()
	s.logger.DebugContext(ctx, "creating card", "id", c.ID)

	now := s.dateTimeTool.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	userID := s.ctxTool.GetReqSubjectFromContext(ctx)
	c.CreatedBy = userID
	c.UpdatedBy = userID

	// store card
	if err := s.repo.CreateCard(ctx, c); err != nil {
		return Card{}, err
	}

	return c, nil
}
