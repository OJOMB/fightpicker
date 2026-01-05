package users

import (
	"context"

	"github.com/pkg/errors"
)

type EmailByTokenVerifier interface {
	VerifyEmailByToken(ctx context.Context, hashedToken []byte) error
}

func (s *Service) VerifyEmailByToken(ctx context.Context, token string) error {
	if token == "" {
		return errors.Wrap(ErrMissingParameter, "token")
	}

	hashedToken, err := s.authTool.HashVerificationToken(token)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to hash verification token", "error", err)
		return ErrInternalError
	}

	if err = s.repo.VerifyEmailByToken(ctx, hashedToken); err != nil {
		return err
	}

	return nil
}
