package fighters

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/id"
)

// FighterLister defines the interface for listing fighters.
type FighterLister interface {
	ListFighters(ctx context.Context, pageSize int, lastSeenID *id.UUID7) ([]Fighter, int, error)
}

// ListFighters retrieves a paginated list of fighters.
// PI is removed from each fighter if the requestor is not an admin.
func (s *Service) ListFighters(ctx context.Context, pageSize int, lastSeenID *id.UUID7) ([]Fighter, int, error) {
	fighters, totalCount, err := s.repo.ListFighters(ctx, pageSize, lastSeenID)
	if err != nil {
		return nil, 0, err
	}

	// we do not inject presigned URLs for list operations to reduce overhead

	return fighters, totalCount, nil
}
