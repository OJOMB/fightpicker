package users

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authservice "github.com/OJOMB/fightpicker/internal/service/auth"
)

func (s *Server) toStatus(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, authservice.ErrMissingParameter):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, authservice.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, authservice.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, authservice.ErrRefreshTokenRevoked):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, authservice.ErrRefreshTokenExpired):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, authservice.ErrRefreshTokenReused):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		s.logger.ErrorContext(ctx, "unexpected internal error", err)
		return status.Error(codes.Internal, "internal error")
	}
}
