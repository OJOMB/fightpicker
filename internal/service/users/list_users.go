package users

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// UserLister defines the interface for listing users.
type UserLister interface {
	ListUsers(ctx context.Context, pageSize int, lastSeenID *uuid.UUID) ([]User, error)
}

// ListUsers retrieves a paginated list of users.
// PI is removed from each user if the requestor is not an admin.
func (svc *Service) ListUsers(ctx context.Context, pageSize int, lastSeenID *uuid.UUID) ([]User, error) {
	users, err := svc.repo.ListUsers(ctx, pageSize, lastSeenID)
	if err != nil {
		return nil, err
	}

	if !contextual.ReqSubjectIsAnAdmin(ctx) {
		for i := range users {
			users[i].removePI()
		}
	}

	// we do not inject presigned URLs for list operations to reduce overhead

	return users, nil
}
