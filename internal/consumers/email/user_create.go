package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UsersService interface {
	SendVerificationEmail(ctx context.Context, userID uuid.UUID, firstName, email string) error
}

type UserCreationEmailConsumer struct {
	client  *kgo.Client
	service UsersService
	logger  logs.Logger
}

func NewUserCreationEmailConsumer(client *kgo.Client, service UsersService, logger logs.Logger) (*UserCreationEmailConsumer, error) {
	if client == nil {
		return nil, ErrNilKafkaClient
	}

	if service == nil {
		return nil, ErrNilUsersService
	}

	if logger == nil {
		return nil, ErrNilLogger
	}

	return &UserCreationEmailConsumer{
		client:  client,
		service: service,
		logger:  logger.With("component", "user_profile_picture_consumer"),
	}, nil
}

func (c *UserCreationEmailConsumer) Run(ctx context.Context) error {
	for {
		fetches := c.client.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return nil
		}

		fetches.EachRecord(func(record *kgo.Record) {
			if err := c.handleRecord(ctx, record); err != nil {
				c.logger.ErrorContext(ctx, "processing failed", "err", err)
				// no commit → retry
				return
			}

			if err := c.client.CommitRecords(ctx, record); err != nil {
				c.logger.ErrorContext(ctx, "commit failed", "err", err)
			}
		})
	}
}

func (c *UserCreationEmailConsumer) handleRecord(ctx context.Context, record *kgo.Record) error {
	userID, err := uuid.FromBytes(record.Key)
	if err != nil {
		return fmt.Errorf("invalid user ID in record key: %w", err)
	}

	var e EventUserCreation
	if err := json.Unmarshal(record.Value, &e); err != nil {
		return fmt.Errorf("unmarshaling record value: %w", err)
	}

	c.logger.InfoContext(ctx, "processing user creation email", "user_id", userID.String())

	if err := c.service.SendVerificationEmail(ctx, userID, e.FirstName, e.Email); err != nil {
		return fmt.Errorf("sending email verification email: %w", err)
	}

	return nil
}
