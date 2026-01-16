package users

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/id"
)

// UserLister defines the interface for listing users.
type UserLister interface {
	ListUsers(ctx context.Context, pageSize int, lastSeenID *id.UUID7) ([]User, int, error)
}

// ListUsers retrieves a paginated list of users.
// PI is removed from each user if the requestor is not an admin.
func (s *Service) ListUsers(ctx context.Context, pageSize int, lastSeenID *id.UUID7) ([]User, int, error) {
	users, totalCount, err := s.repo.ListUsers(ctx, pageSize, lastSeenID)
	if err != nil {
		return nil, 0, err
	}

	if !s.ctxTool.ReqSubjectIsAnAdmin(ctx) {
		for i := range users {
			users[i].removePI()
		}
	}

	// we do not inject presigned URLs for list operations to reduce overhead

	return users, totalCount, nil
}
