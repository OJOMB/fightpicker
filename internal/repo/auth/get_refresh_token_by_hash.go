package auth

import (
	"context"

	authservice "github.com/OJOMB/fightpicker/internal/service/auth"
)

// GetRefreshTokenByHash retrieves a refresh token by its hash.
func (r *Repo) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (authservice.RefreshToken, error) {
	row, err := r.client.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return authservice.RefreshToken{}, dbErrorToRepoError(err)
	}

	return refreshTokenDBOtoIDO(row), nil
}
