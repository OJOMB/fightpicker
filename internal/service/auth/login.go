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
		s.logger.DebugContext(ctx, "user not found with email: "+email)
		// for security purposes, we do not explicitly reveal whether the email is registered or not
		return "", "", time.Time{}, ErrInvalidCredentials
	}

	if ok, err := s.passwordVerifier.Verify(user.PasswordHash, password); err != nil {
		return "", "", time.Time{}, err
	} else if !ok {
		s.logger.DebugContext(ctx, "user password incorrect: "+email)
		// we don't let the user know it's specifically the password that is incorrect for security reasons
		// user enumeration protection: we return a generic invalid credentials error instead of specifically revealing that the password is incorrect
		// thus not revealing whether the email is registered or not
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
		map[string]any{"roles": roles, "perms": permissions, "type": "access"},
		s.secretKey,
	)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// generate Refresh JWT token
	refreshToken, err := s.jwtGen.GenerateToken(
		s.id.Generate(),
		user.ID,
		s.refreshTokenTTL,
		s.tokenIssuer,
		s.tokenAudience,
		map[string]any{"type": "refresh", "family_id": s.id.Generate()},
		s.secretKey,
	)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// store the refresh token in the database
	refreshTokenHash := s.jwtGen.HashTokenString(refreshToken.TokenStr)

	// extract IP address and User-Agent from the context for auditing and security purposes
	// in both cases, we don't fail the login process if we can't determine legitimate values
	// these values are auditing metadata (for “active sessions” and spotting suspicious activity), not strictly part of authenticating the user.
	// empty IP address should never really happen irl, empty User-Agent is also uncommon but possible (e.g. some bots or privacy-focused browsers)
	ipAddress, ok := ctx.Value(contextual.KeyRequestRemoteIPAddr).(string)
	if !ok || ipAddress == "" {
		s.logger.WarnContext(ctx, "could not determine IP address for login")
	}

	userAgent, ok := ctx.Value(contextual.KeyRequestUserAgent).(string)
	if !ok || userAgent == "" {
		s.logger.WarnContext(ctx, "could not determine User-Agent for login")
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
