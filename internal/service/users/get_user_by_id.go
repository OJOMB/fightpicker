package users

import (
	"context"
	"net/http"

	"github.com/OJOMB/fightpicker/pkg/id"
)

type PresignedGetURLGenerator interface {
	GeneratePresignedGetURL(ctx context.Context, userID id.UUID7) (string, http.Header, error)
}

// UserByIDGetter defines the interface for retrieving a user by ID.
type UserByIDGetter interface {
	GetUserByID(ctx context.Context, id id.UUID7) (User, error)
}

// GetUserByID retrieves a user by their uuid.
func (s *Service) GetUserByID(ctx context.Context, id id.UUID7) (User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return User{}, err
	}

	// if the requesting user is not the user being requested and is not an admin, remove sensitive fields
	if s.ctxTool.GetReqSubjectFromContext(ctx) != id && !s.ctxTool.ReqSubjectIsAnAdmin(ctx) {
		user.removePI()
	}

	s.injectPresignedGetURLIfNeeded(ctx, &user)

	return user, nil
}
