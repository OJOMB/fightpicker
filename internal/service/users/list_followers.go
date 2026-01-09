package users

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

type UserFollowerLister interface {
	ListFollowers(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]User, int, error)
}

// ListFollowers retrieves a paginated list of followers for a given user
// PI is removed from each follower if the requestor is not an admin.
func (s *Service) ListFollowers(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]User, int, error) {
	followers, totalCount, err := s.repo.ListFollowers(ctx, userID, pageSize, lastSeenID)
	if err != nil {
		return nil, 0, err
	}

	if !contextual.ReqSubjectIsAnAdmin(ctx) {
		for i := range followers {
			followers[i].removePI()
		}
	}

	// we do not inject presigned URLs for list operations to reduce overhead

	return followers, totalCount, nil
}
