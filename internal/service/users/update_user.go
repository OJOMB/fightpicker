package users

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// UserUpdater defines the interface for updating users.
type UserUpdater interface {
	UpdateUser(ctx context.Context, id uuid.UUID, updates UserUpdate) (User, error)
}

// UpdateUser updates the user with the given ID using the provided updates.
func (svc *Service) UpdateUser(ctx context.Context, id uuid.UUID, updates UserUpdate) (User, error) {
	reqSubject := contextual.GetReqSubjectFromContext(ctx)
	if reqSubject == uuid.Nil {
		svc.logger.ErrorContext(ctx, "unable to get request subject from context")
		return User{}, ErrInternalError
	}

	if !contextual.ReqSubjectIsAnAdmin(ctx) && reqSubject != id {
		return User{}, ErrUnauthorized
	}

	if err := svc.validateUpdateReq(updates); err != nil {
		return User{}, err
	}

	updates.UpdatedAt = svc.dateTimeTool.Now()
	updates.UpdatedBy = reqSubject

	if updates.Password != nil {
		hashedPassword, err := svc.authTool.HashPassword(*updates.Password)
		if err != nil {
			svc.logger.ErrorContext(ctx, "failed to hash password", "error", err, "user_id", id)
			return User{}, err
		}

		updates.Password = &hashedPassword
	}

	user, err := svc.repo.UpdateUser(ctx, id, updates)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
