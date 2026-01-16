package auth

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// Login authenticates a user using their email and password, returning an access, refresh token and its expiration time.
func (s *Service) Login(ctx context.Context, email, password string) (string, string, time.Time, error) {
	if email == "" {
		return "", "", time.Time{}, errors.Wrap(ErrMissingParameter, "email")
	}

	if password == "" {
		return "", "", time.Time{}, errors.Wrap(ErrMissingParameter, "password")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// TODO: distinguish between not found and other errors
		return "", "", time.Time{}, ErrUserNotFound
	}

	if ok, err := s.passwordVerifier.Verify(user.PasswordHash, password); err != nil {
		return "", "", time.Time{}, err
	} else if !ok {
		return "", "", time.Time{}, ErrInvalidCredentials
	}

	roles, permissions, err := s.authRepo.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// generate Access JWT token
	accessToken, err := s.jwtGen.GenerateToken(
		s.id.Generate(),
		user.ID,
		s.accessTokenTTL,
		s.tokenIssuer,
		s.tokenAudience,
		map[string]any{
			"roles": roles, "perms": permissions, "type": "access",
		},
		s.secretKey,
	)

	// generate Refresh JWT token
	refreshToken, err := s.jwtGen.GenerateToken(
		s.id.Generate(),
		user.ID,
		s.refreshTokenTTL,
		s.tokenIssuer,
		s.tokenAudience,
		map[string]any{"type": "refresh"},
		s.secretKey,
	)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// store the refresh token in the database
	refreshTokenHash := s.jwtGen.HashTokenString(refreshToken.TokenStr)
	ipAddress, ok := ctx.Value(contextual.KeyRequestRemoteAddr).(string)
	if !ok {
		s.logger.WarnContext(ctx, "could not get IP address from context")
	}

	userAgent, ok := ctx.Value(contextual.KeyRequestUserAgent).(string)
	if !ok {
		s.logger.WarnContext(ctx, "could not get User-Agent from context")
	}

	if err := s.authRepo.StoreRefreshToken(ctx, user.ID, refreshToken.JTI, refreshTokenHash, ipAddress, userAgent, refreshToken.ExpiresAt); err != nil {
		return "", "", time.Time{}, err
	}

	// finally update last login time
	if err := s.userRepo.UpdateLastLoginAtByUserID(ctx, user.ID); err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken.TokenStr, refreshToken.TokenStr, refreshToken.ExpiresAt, nil
}
