package users

import (
	"context"

	"github.com/pkg/errors"
)

// UserByUsernameGetter defines the interface for retrieving a user by username.
type UserByUsernameGetter interface {
	GetUserByUsername(ctx context.Context, username string) (User, error)
}

// GetUserByUsername retrieves a user by their username.
func (s *Service) GetUserByUsername(ctx context.Context, username string) (User, error) {
	if username == "" {
		return User{}, errors.Wrap(ErrMissingParameter, "username")
	}

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return User{}, err
	}

	reqSubjectID := s.ctxTool.GetReqSubjectFromContext(ctx)
	if reqSubjectID != user.ID && !s.ctxTool.ReqSubjectIsAnAdmin(ctx) {
		user.removePI()
	}

	s.injectPresignedGetURLIfNeeded(ctx, &user)

	return user, nil
}
