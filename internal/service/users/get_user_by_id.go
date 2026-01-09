package users

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

type PresignedGetURLGenerator interface {
	GeneratePresignedGetURL(ctx context.Context, userID uuid.UUID) (string, http.Header, error)
}

// UserByIDGetter defines the interface for retrieving a user by ID.
type UserByIDGetter interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
}

// GetUserByID retrieves a user by their uuid.
func (svc *Service) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := svc.repo.GetUserByID(ctx, id)
	if err != nil {
		return User{}, err
	}

	// if the requesting user is not the user being requested and is not an admin, remove sensitive fields
	if contextual.GetReqSubjectFromContext(ctx) != id && !contextual.ReqSubjectIsAnAdmin(ctx) {
		user.removePI()
	}

	svc.injectPresignedGetURLIfNeeded(ctx, &user)

	return user, nil
}
