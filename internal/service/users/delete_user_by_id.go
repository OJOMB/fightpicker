package users

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/id"
)

// UserByIDGetter defines the interface for retrieving a user by ID.
type UserByIDDeleter interface {
	DeleteUserByID(ctx context.Context, id id.UUID7) error
}

// DeleteUserByID deletes a user by their uuid.
func (svc *Service) DeleteUserByID(ctx context.Context, id id.UUID7) error {
	reqSubject := svc.ctxTool.GetReqSubjectFromContext(ctx)
	if reqSubject != id && !svc.ctxTool.ReqSubjectIsAnAdmin(ctx) {
		return ErrUnauthorized
	}

	err := svc.repo.DeleteUserByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
