package auth

import (
	"context"

	service "github.com/OJOMB/fightpicker/internal/service/auth"
)

// GetRefreshTokenByHash retrieves a refresh token by its hash.
func (r *Repo) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (service.RefreshToken, error) {
	row, err := r.client.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return service.RefreshToken{}, dbErrorToServiceError(err)
	}

	return refreshTokenDBOtoIDO(row), nil
}
