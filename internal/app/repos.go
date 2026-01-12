package app

import (
	"context"
	"time"

	"github.com/OJOMB/fightpicker/internal/config"
	authrepo "github.com/OJOMB/fightpicker/internal/repo/auth"
	fightersrepo "github.com/OJOMB/fightpicker/internal/repo/fighters"
	usersrepo "github.com/OJOMB/fightpicker/internal/repo/users"
)

type repos struct {
	AuthRepo     *authrepo.Repo
	UsersRepo    *usersrepo.Repo
	FightersRepo *fightersrepo.Repo
}

func (a *App) newRepos(ctx context.Context, cfg *config.Config) error {
	authRepo := authrepo.New(a.DB.pool, a.DB.queries, a.Utils.IDTool, a.Logger)

	presignedGetURLTTLDuration := time.Duration(cfg.AWS.PresignedGetURLTTL) * time.Minute
	presignedPutURLTTLDuration := time.Duration(cfg.AWS.PresignedPutURLTTL) * time.Minute

	usersRepo, err := usersrepo.New(
		a.DB.pool,
		cfg.HTTP.Domain,
		a.DB.queries,
		a.Clients.S3Client,
		a.Utils.DateTimeTool,
		cfg.AWS.S3MediaBucket,
		cfg.AWS.S3MediaBucketRaw,
		presignedGetURLTTLDuration,
		presignedPutURLTTLDuration,
		a.Clients.KafkaPostUserCreateClient,
		cfg.EventBroker.TopicPostUserCreated,
		cfg.EventBroker.TopicPostUserVerified,
		a.Clients.EmailDialer,
		cfg.Email.AddressNoReply,
		a.Logger,
	)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create users repo", "error", err)

	}

	fightersrepo, err := fightersrepo.New(a.DB.pool, a.DB.queries, a.Clients.RedisClient, a.Utils.DateTimeTool, a.Logger)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create fighters repo", "error", err)
		return err
	}

	a.Repos = repos{
		AuthRepo:     authRepo,
		UsersRepo:    usersRepo,
		FightersRepo: fightersrepo,
	}

	return nil
}
