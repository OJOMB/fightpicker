package users

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/id"
)

// UserUpdater defines the interface for updating users.
type UserUpdater interface {
	UpdateUser(ctx context.Context, id id.UUID7, updates UserUpdate) (User, error)
}

// UpdateUser updates the user with the given ID using the provided updates.
func (s *Service) UpdateUser(ctx context.Context, userID id.UUID7, updates UserUpdate) (User, error) {
	reqSubject := s.ctxTool.GetReqSubjectFromContext(ctx)
	if reqSubject == id.UUID7Nil {
		s.logger.ErrorContext(ctx, "unable to get request subject from context")
		return User{}, ErrInternalError
	}

	if !s.ctxTool.ReqSubjectIsAnAdmin(ctx) && reqSubject != userID {
		return User{}, ErrUnauthorized
	}

	if err := s.validateUpdateReq(updates); err != nil {
		return User{}, err
	}

	updates.UpdatedAt = s.dateTimeTool.Now()
	updates.UpdatedBy = reqSubject

	if updates.Password != nil {
		hashedPassword, err := s.authTool.HashPassword(*updates.Password)
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to hash password", "error", err, "user_id", userID)
			return User{}, err
		}

		updates.Password = &hashedPassword
	}

	user, err := s.repo.UpdateUser(ctx, userID, updates)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
