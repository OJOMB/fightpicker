package users

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/id"
)

type UserFolloweeLister interface {
	ListFollowees(ctx context.Context, userID id.UUID7, pageSize int, lastSeenID *id.UUID7) ([]User, int, error)
}

// ListFollowees retrieves a paginated list of followees for a given user
// PI is removed from each followee if the requestor is not an admin.
func (s *Service) ListFollowees(ctx context.Context, userID id.UUID7, pageSize int, lastSeenID *id.UUID7) ([]User, int, error) {
	followees, totalCount, err := s.repo.ListFollowees(ctx, userID, pageSize, lastSeenID)
	if err != nil {
		return nil, 0, err
	}

	if !s.ctxTool.ReqSubjectIsAnAdmin(ctx) {
		for i := range followees {
			followees[i].removePI()
		}
	}

	// we do not inject presigned URLs for list operations to reduce overhead

	return followees, totalCount, nil
}
