package users

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// UserByEmailGetter defines the interface for retrieving a user by email.
type UserByEmailGetter interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

// GetUserByEmail retrieves a user by their email address.
func (svc *Service) GetUserByEmail(ctx context.Context, email string) (User, error) {
	user, err := svc.repo.GetUserByEmail(ctx, email)
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
