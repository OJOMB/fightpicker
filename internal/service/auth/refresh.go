package auth

import (
	"context"
	"time"

	"github.com/OJOMB/fightpicker/pkg/contextual"
)

// Refresh validates a refresh token and generates new access and refresh tokens.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, string, time.Time, error) {
	hashedToken := s.jwtGen.HashTokenString(refreshToken)

	currentRefreshToken, err := s.authRepo.GetRefreshTokenByHash(ctx, hashedToken)
	if err != nil {
		s.logger.DebugContext(ctx, "failed to get refresh token by hash", "error", err)
		return "", "", time.Time{}, err
	}

	// check the validity of the current refresh token
	if currentRefreshToken.Revoked {
		s.logger.DebugContext(ctx, "refresh token has been revoked", "token_id", currentRefreshToken.ID)
		return "", "", time.Time{}, ErrRefreshTokenRevoked
	}

	if currentRefreshToken.ExpiresAt.Before(time.Now()) {
		s.logger.DebugContext(ctx, "refresh token has expired", "token_id", currentRefreshToken.ID)
		return "", "", time.Time{}, ErrRefreshTokenExpired
	}

	if currentRefreshToken.ReplacedBy != nil {
		s.logger.DebugContext(ctx, "refresh token has been replaced, possible reuse detected", "token_id", currentRefreshToken.ID)
		return "", "", time.Time{}, ErrRefreshTokenReused
	}

	userID := currentRefreshToken.UserID

	roles, permissions, err := s.authRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// generate new Access JWT token
	accessToken, err := s.jwtGen.GenerateToken(
		s.idGen.Generate(),
		userID,
		s.accessTokenTTL,
		s.tokenIssuer,
		s.tokenAudience,
		map[string]any{
			"roles": roles, "perms": permissions, "type": tokenTypeAccess,
		},
		s.secretKey,
	)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// generate new Refresh JWT token
	refreshTokenNew, err := s.jwtGen.GenerateToken(
		s.idGen.Generate(),
		userID,
		s.refreshTokenTTL,
		s.tokenIssuer,
		s.tokenAudience,
		map[string]any{"type": tokenTypeRefresh},
		s.secretKey,
	)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// store the refresh token in the database
	refreshTokenHash := s.jwtGen.HashTokenString(refreshTokenNew.TokenStr)
	ipAddress, ok := ctx.Value(contextual.KeyRequestRemoteAddr).(string)
	if !ok {
		ipAddress = ""
	}

	userAgent, ok := ctx.Value(contextual.KeyRequestUserAgent).(string)
	if !ok {
		s.logger.WarnContext(ctx, "failed to get user agent from context")
		userAgent = ""
	}

	if err := s.authRepo.RotateRefreshTokens(ctx, hashedToken, refreshTokenHash, refreshTokenNew.JTI, userID, ipAddress, userAgent, refreshTokenNew.ExpiresAt); err != nil {
		s.logger.ErrorContext(ctx, "failed to rotate refresh tokens", "error", err, "token_id", currentRefreshToken.ID)
		return "", "", time.Time{}, err
	}

	return accessToken.TokenStr, refreshTokenNew.TokenStr, refreshTokenNew.ExpiresAt, nil
}
