package users

import (
	"context"

	"github.com/pkg/errors"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// UserByUsernameGetter defines the interface for retrieving a user by username.
type UserByUsernameGetter interface {
	GetUserByUsername(ctx context.Context, username string) (User, error)
}

// GetUserByUsername retrieves a user by their username.
func (svc *Service) GetUserByUsername(ctx context.Context, username string) (User, error) {
	if username == "" {
		return User{}, errors.Wrap(ErrMissingParameter, "username")
	}

	user, err := svc.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return User{}, err
	}

	reqSubjectID := contextual.GetReqSubjectFromContext(ctx)
	if reqSubjectID != user.ID && !contextual.ReqSubjectIsAnAdmin(ctx) {
		user.removePI()
	}

	svc.injectPresignedGetURLIfNeeded(ctx, &user)

	return user, nil
}
