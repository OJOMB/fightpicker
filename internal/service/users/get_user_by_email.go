package users

import (
	"context"
)

// UserByEmailGetter defines the interface for retrieving a user by email.
type UserByEmailGetter interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

// GetUserByEmail retrieves a user by their email address.
func (s *Service) GetUserByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
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
