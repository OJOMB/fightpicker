package auth

import (
	"context"
	"time"
)

// Logout invalidates a refresh token, effectively logging the user out.
// New access tokens cannot be generated using the invalidated refresh token.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	// hash the provided refresh token
	hashedToken := s.jwtGen.HashTokenString(refreshToken)

	// TODO: it's probably unnecessary to check the token's validity first - we can just attempt to revoke it in the repo
	// leaving it for now

	// get the refresh token row from the database
	currentRefreshToken, err := s.authRepo.GetRefreshTokenByHash(ctx, hashedToken)
	if err != nil {
		return err
	}

	if currentRefreshToken.Revoked || currentRefreshToken.ExpiresAt.Before(time.Now()) {
		// Token is already invalid - updating it anyway is unnecessary since it adds no security benefit
		return nil
	}

	if err := s.authRepo.RevokeRefreshTokenByHash(ctx, hashedToken); err != nil {
		return err
	}

	return nil
}
