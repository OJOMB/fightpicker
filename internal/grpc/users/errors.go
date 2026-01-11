package users

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	usersrepo "github.com/OJOMB/fightpicker/internal/repo/users"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
)

func (s *Server) toStatus(err error) error {
	switch {
	case errors.Is(err, usersrepo.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, usersrepo.ErrEmailTaken),
		errors.Is(err, usersrepo.ErrUsernameTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, usersservice.ErrInvalidParameter),
		errors.Is(err, usersservice.ErrMissingParameter):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, usersservice.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, usersrepo.ErrDefaultRoleNotFound):
		return status.Error(codes.FailedPrecondition, "system setup error")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
