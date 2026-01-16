package users

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/id"
)

// UserFollower defines the interface for following users.
type UserFollower interface {
	FollowUser(ctx context.Context, followID, followerID, followeeID id.UUID7) error
}

// FollowUser allows a user to follow another user by their ID.
func (svc *Service) FollowUser(ctx context.Context, followeeID id.UUID7) error {
	reqSubject := svc.ctxTool.GetReqSubjectFromContext(ctx)

	followID := svc.idTool.Generate()

	if err := svc.repo.FollowUser(ctx, followID, reqSubject, followeeID); err != nil {
		return err
	}

	return nil
}
