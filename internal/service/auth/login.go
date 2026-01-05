package auth

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// Login authenticates a user using their email and password, returning an access and refresh token.
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

	if ok, err := s.pwordVerifier.Verify(user.PasswordHash, password); err != nil {
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
		s.idGen.Generate(),
		user.ID,
		s.accessTokenTTL,
		"fightpicker",
		"fightpicker_users",
		map[string]any{
			"roles": roles, "perms": permissions, "type": "access",
		},
		s.secretKey,
	)

	// generate Refresh JWT token
	refreshToken, err := s.jwtGen.GenerateToken(
		s.idGen.Generate(),
		user.ID,
		s.refreshTokenTTL,
		"fightpicker",
		"fightpicker_users",
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
		ipAddress = ""
	}

	userAgent, ok := ctx.Value(contextual.KeyRequestUserAgent).(string)
	if !ok {
		userAgent = ""
	}

	if err := s.authRepo.StoreRefreshToken(ctx, user.ID, refreshToken.JTI, refreshTokenHash, ipAddress, userAgent, refreshToken.ExpiresAt); err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken.TokenStr, refreshToken.TokenStr, refreshToken.ExpiresAt, nil
}
