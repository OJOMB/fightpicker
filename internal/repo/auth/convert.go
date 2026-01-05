package auth

import (
	service "github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

func refreshTokenDBOtoIDO(dbToken postgres.RefreshToken) service.RefreshToken {
	return service.RefreshToken{
		ID:        dbToken.ID,
		UserID:    dbToken.UserID,
		TokenHash: dbToken.TokenHash,
		JTI:       dbToken.Jti,
		ExpiresAt: dbToken.ExpiresAt.Time,
		Revoked:   dbToken.Revoked,
		IPAddress: dbToken.IpAddress.String,
		UserAgent: dbToken.UserAgent.String,
		CreatedAt: dbToken.CreatedAt.Time,
		UpdatedAt: dbToken.UpdatedAt.Time,
	}
}
