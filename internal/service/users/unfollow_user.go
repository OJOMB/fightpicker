package users

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// UserUnfollower defines the interface for unfollowing users.
type UserUnfollower interface {
	UnfollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error
}

// UnfollowUser allows a user to unfollow another user by their ID.
func (svc *Service) UnfollowUser(ctx context.Context, followeeID uuid.UUID) error {
	reqSubject := contextual.GetReqSubjectFromContext(ctx)
	if reqSubject == uuid.Nil {
		svc.logger.ErrorContext(ctx, "unable to get request subject from context")
		return ErrInternalError
	}

	if err := svc.repo.UnfollowUser(ctx, reqSubject, followeeID); err != nil {
		return err
	}

	return nil
}
