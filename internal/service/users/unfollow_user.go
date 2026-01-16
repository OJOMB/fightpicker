package users

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/id"
)

// UserUnfollower defines the interface for unfollowing users.
type UserUnfollower interface {
	UnfollowUser(ctx context.Context, followerID, followeeID id.UUID7) error
}

// UnfollowUser allows a user to unfollow another user by their ID.
func (s *Service) UnfollowUser(ctx context.Context, followeeID id.UUID7) error {
	reqSubject := s.ctxTool.GetReqSubjectFromContext(ctx)
	if reqSubject == id.UUID7Nil {
		s.logger.ErrorContext(ctx, "unable to get request subject from context")
		return ErrInternalError
	}

	if err := s.repo.UnfollowUser(ctx, reqSubject, followeeID); err != nil {
		return err
	}

	return nil
}
