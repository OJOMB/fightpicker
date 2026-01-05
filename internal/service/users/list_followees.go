package users

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

type UserFolloweeLister interface {
	ListFollowees(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]User, error)
}

// ListFollowees retrieves a paginated list of followees for a given user
// PI is removed from each followee if the requestor is not an admin.
func (s *Service) ListFollowees(ctx context.Context, userID uuid.UUID, pageSize int, lastSeenID *uuid.UUID) ([]User, error) {
	followees, err := s.repo.ListFollowees(ctx, userID, pageSize, lastSeenID)
	if err != nil {
		return nil, err
	}

	if !contextual.ReqSubjectIsAnAdmin(ctx) {
		for i := range followees {
			followees[i].removePI()
		}
	}

	// we do not inject presigned URLs for list operations to reduce overhead

	return followees, nil
}
