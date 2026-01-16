package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UsersService interface {
	SendVerificationEmail(ctx context.Context, userID id.UUID7, firstName, email string) error
}

type UserCreationEmailConsumer struct {
	client  *kgo.Client
	service UsersService
	id      id.UUID7Parser
	logger  logs.Logger
}

func NewUserCreationEmailConsumer(client *kgo.Client, service UsersService, idTool id.UUID7Parser, logger logs.Logger) (*UserCreationEmailConsumer, error) {
	if client == nil {
		return nil, ErrNilKafkaClient
	}

	if service == nil {
		return nil, ErrNilUsersService
	}

	if idTool == nil {
		return nil, ErrNilIDTool
	}

	if logger == nil {
		return nil, ErrNilLogger
	}

	return &UserCreationEmailConsumer{
		client:  client,
		service: service,
		id:      idTool,
		logger:  logger.With("component", "user_profile_picture_consumer"),
	}, nil
}

// Run starts the consumer loop to process incoming Kafka records.
func (c *UserCreationEmailConsumer) Run(ctx context.Context) error {
	c.logger.InfoContext(ctx, "user creation email consumer is running")

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
	userID, err := c.id.ParseBytes(record.Key)
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
