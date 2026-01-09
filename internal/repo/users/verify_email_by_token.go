package users

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/OJOMB/fightpicker/pkg/clients/postgres"
)

// VerifyEmailByToken verifies a user's email using the provided hashed token.
func (r *Repo) VerifyEmailByToken(ctx context.Context, hashedToken []byte) error {
	r.metrics.EmailVerificationAttempts.Add(ctx, 1)

	userID, err := r.dbClient.VerifyUserEmailByTokenHash(ctx, postgres.VerifyUserEmailByTokenHashParams{
		EmailVerificationTokenHash: hashedToken,
		UpdatedAt: pgtype.Timestamptz{
			Time:  r.dateTimeTool.Now().UTC(),
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	if userID == uuid.Nil {
		r.logger.WarnContext(ctx, "email verification token is invalid or expired", "user_id", userID, "hashed_token", hashedToken)
		r.metrics.EmailVerificationFailures.Add(ctx, 1)
		return nil
	}

	// now publish an event to Kafka indicating the user's email has been verified and welcome email can be sent
	r.publishUserVerified(userID)

	return nil
}

func (r *Repo) publishUserVerified(userID uuid.UUID) {
	if r.events.client == nil {
		return
	}

	key := userID.Bytes()

	// use a fresh context - we don't want the req ctx to time out because the user creation request has finished before the callback is invoked
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	r.events.client.Produce(ctx, &kgo.Record{
		Topic: r.events.topicPostUserVerified,
		Key:   key,
	}, func(record *kgo.Record, err error) {
		defer cancel()
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to produce user.verified", "error", err, "record_key", record.Key)
		}
	})
}
