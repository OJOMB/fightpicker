package users

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// UserFollower defines the interface for following users.
type UserFollower interface {
	FollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error
}

// FollowUser allows a user to follow another user by their ID.
func (svc *Service) FollowUser(ctx context.Context, followeeID uuid.UUID) error {
	reqSubject := contextual.GetReqSubjectFromContext(ctx)
	if reqSubject == uuid.Nil {
		svc.logger.ErrorContext(ctx, "unable to get request subject from context")
		return ErrInternalError
	}

	if err := svc.repo.FollowUser(ctx, reqSubject, followeeID); err != nil {
		return err
	}

	return nil
}
