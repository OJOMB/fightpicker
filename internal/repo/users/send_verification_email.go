package users

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"gopkg.in/mail.v2"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/id"
)

// SendVerificationEmail generates a verification token, stores its hash in the database, and sends a verification email to the user.
func (r *Repo) SendVerificationEmail(ctx context.Context, userID id.UUID7, hashedToken []byte, msg *mail.Message) error {
	now := r.dateTimeTool.Now().UTC()

	// store the token in the database associated with the user
	if err := r.dbClient.UpdateEmailVerificationTokenHashByUserID(ctx, postgres.UpdateEmailVerificationTokenHashByUserIDParams{
		ID:                         userID,
		EmailVerificationTokenHash: hashedToken,
		EmailVerificationTokenExpiresAt: pgtype.Timestamptz{
			Time:  now.Add(1 * time.Hour),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("failed to store verification token: %w", err)
	}

	if err := r.email.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send verification email to: %w", err)
	}

	return nil
}
