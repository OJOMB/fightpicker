package auth

import (
	"context"

	service "github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

func (r *Repo) refreshTokenDBOtoIDO(ctx context.Context, dbToken postgres.RefreshToken) service.RefreshToken {
	var replacedBy *id.UUID7
	if dbToken.ReplacedBy.Valid {
		val, err := r.id.ParseString(dbToken.ReplacedBy.String())
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to parse refresh token replacedBy field", err)
		} else {
			replacedBy = &val
		}
	}

	return service.RefreshToken{
		ID:         dbToken.ID,
		FamilyID:   dbToken.FamilyID,
		UserID:     dbToken.UserID,
		TokenHash:  dbToken.TokenHash,
		JTI:        dbToken.Jti,
		ExpiresAt:  dbToken.ExpiresAt.Time,
		Revoked:    dbToken.Revoked,
		ReplacedBy: replacedBy,
		IPAddress:  dbToken.IpAddress.String,
		UserAgent:  dbToken.UserAgent.String,
		CreatedAt:  dbToken.CreatedAt.Time,
		UpdatedAt:  dbToken.UpdatedAt.Time,
	}
}
