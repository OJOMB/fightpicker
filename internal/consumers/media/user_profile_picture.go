package media

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/OJOMB/fightpicker/pkg/logs"
)

type UsersService interface {
	ProcessUploadedUserProfilePicture(ctx context.Context, userID uuid.UUID, bucketName, objectKey string) error
}

type UserProfilePictureConsumer struct {
	client  *kgo.Client
	service UsersService
	logger  logs.Logger
}

func NewUserProfilePictureConsumer(client *kgo.Client, service UsersService, logger logs.Logger) (*UserProfilePictureConsumer, error) {
	if client == nil {
		return nil, ErrNilKafkaClient
	}

	if service == nil {
		return nil, ErrNilUsersService
	}

	if logger == nil {
		return nil, ErrNilLogger
	}

	return &UserProfilePictureConsumer{
		client:  client,
		service: service,
		logger:  logger.With("component", "user_profile_picture_consumer"),
	}, nil
}

func (c *UserProfilePictureConsumer) Run(ctx context.Context) error {
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

func (c *UserProfilePictureConsumer) handleRecord(ctx context.Context, record *kgo.Record) error {
	// the record key is the object file path, which includes the user ID
	// e.g. fightpicker-media-raw/users/019b38cb-8cf2-7ec5-897f-005b983dbd00/media/profile_image.jpeg
	pathComponents := strings.Split(string(record.Key), "/")
	if len(pathComponents) < 3 {
		c.logger.ErrorContext(ctx, "invalid record key format", "key", string(record.Key))
		return ErrInvalidRecordKeyFormat
	}

	userIDStr := pathComponents[2]
	userUUID, err := uuid.FromString(userIDStr)
	if err != nil {
		c.logger.ErrorContext(ctx, "invalid user ID in record key", "user_id", userIDStr, "err", err)
		return ErrInvalidRecordKeyFormat
	}

	// unmarshal the record value to get the s3event
	var event EventAWS
	if err := json.Unmarshal(record.Value, &event); err != nil {
		c.logger.ErrorContext(ctx, "failed to unmarshal s3 event", "err", err, "value", record.Value)
		return err
	}

	c.logger.DebugContext(ctx, "processing %s from source %s and region %s", event.EventName, event.EventSource, event.AWSRegion)

	for _, record := range event.Records {
		bucketName := record.S3.Bucket.Name
		objectKey, err := url.PathUnescape(record.S3.Object.Key)
		if err != nil {
			c.logger.ErrorContext(ctx, "failed to unescape object key", "error", err)
			return err
		}

		c.logger.InfoContext(ctx, "processing uploaded user profile picture", "user_id", userIDStr, "bucket", bucketName, "object_key", objectKey)

		if err := c.service.ProcessUploadedUserProfilePicture(ctx, userUUID, bucketName, objectKey); err != nil {
			c.logger.ErrorContext(ctx, "failed to process uploaded user profile picture", "err", err, "user_id", userIDStr, "bucket", bucketName, "object_key", objectKey)
			return err
		}
	}

	return nil
}
