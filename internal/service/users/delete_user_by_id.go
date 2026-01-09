package users

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// UserByIDGetter defines the interface for retrieving a user by ID.
type UserByIDDeleter interface {
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
}

// DeleteUserByID deletes a user by their uuid.
func (svc *Service) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	reqSubject := contextual.GetReqSubjectFromContext(ctx)
	if reqSubject != id && !contextual.ReqSubjectIsAnAdmin(ctx) {
		return ErrUnauthorized
	}

	err := svc.repo.DeleteUserByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
