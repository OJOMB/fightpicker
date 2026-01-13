package app

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"gopkg.in/mail.v2"

	"github.com/OJOMB/fightpicker/internal/config"
)

type clients struct {
	S3Client                   *s3.Client
	RedisClient                *redis.Client
	EmailDialer                *mail.Dialer
	KafkaMediaProfilePicClient *kgo.Client
	KafkaPostUserCreateClient  *kgo.Client
}

func (a *App) newClients(ctx context.Context, cfg *config.Config) error {
	s3Client, err := a.newAWSClients(ctx, cfg.AWS)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create AWS S3 client", "error", err)
		return err
	}

	redisClient, err := a.newRedisClient(ctx, cfg.Cache)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create Redis client", "error", err)
		return err
	}

	emailDialer := a.newEmailDialer(cfg.Email)

	kafkaMediaProfilePicClient, kafkaPostUserCreateClient, err := a.newKafkaClients(ctx, cfg.EventBroker)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create Kafka clients", "error", err)
		return err
	}

	a.Clients = clients{
		S3Client:                   s3Client,
		RedisClient:                redisClient,
		EmailDialer:                emailDialer,
		KafkaMediaProfilePicClient: kafkaMediaProfilePicClient,
		KafkaPostUserCreateClient:  kafkaPostUserCreateClient,
	}

	return nil
}

func (a *App) newAWSClients(ctx context.Context, cfg config.AWSConfig) (*s3.Client, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awscfg, func(o *s3.Options) {
		if cfg.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = true
			creds, err := o.Credentials.Retrieve(ctx)
			if err != nil {
				a.Logger.ErrorContext(ctx, "failed to load AWS credentials", "access_key", creds.AccessKeyID)

				return
			}
		}
	})

	return s3Client, nil
}

func (a *App) newRedisClient(ctx context.Context, cfg config.CacheConfig) (*redis.Client, error) {
	// Redis Client init
	rdbClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0, // use default DB
		Protocol: 2,
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(rdbClient); err != nil {
		a.Logger.FatalContext(ctx, "failed to instrument redis client for tracing", "error", err)
		return nil, err
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(rdbClient); err != nil {
		a.Logger.FatalContext(ctx, "failed to instrument redis client for metrics", "error", err)
		return nil, err
	}

	a.Logger.InfoContext(ctx, "connected to redis cache", "host", cfg.Host, "port", cfg.Port)

	return rdbClient, nil
}

func (a *App) newEmailDialer(cfg config.EmailConfig) *mail.Dialer {
	emailDialer := mail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)
	if cfg.SkipTLS {
		emailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return emailDialer
}

func (a *App) newKafkaClients(ctx context.Context, cfg config.EventBrokerConfig) (*kgo.Client, *kgo.Client, error) {
	kafkaMediaProfilePicClient, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.SeedBrokers...),
		kgo.ConsumeTopics(cfg.TopicProfilePictureUpload),
		// consumer group ID is is required for CommitRecords to work
		kgo.ConsumerGroup(cfg.GroupID),
		// start from the earliest offset if there is no committed offset
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		a.Logger.FatalContext(ctx, "failed to create kafka client for media profile picture uploads", "error", err)
		return nil, nil, err
	}

	kafkaPostUserCreateClient, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.SeedBrokers...),
		kgo.ConsumeTopics(cfg.TopicPostUserCreated),
		// consumer group ID is is required for CommitRecords to work
		kgo.ConsumerGroup(cfg.GroupID),
		// start from the earliest offset if there is no committed offset
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		a.Logger.FatalContext(ctx, "failed to create kafka client for post user create events", "error", err)
	}

	return kafkaMediaProfilePicClient, kafkaPostUserCreateClient, nil
}
