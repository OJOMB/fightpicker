package app

import (
	"context"
	"time"

	"github.com/OJOMB/fightpicker/internal/config"
	authrepo "github.com/OJOMB/fightpicker/internal/repo/auth"
	cardsrepo "github.com/OJOMB/fightpicker/internal/repo/cards"
	fightersrepo "github.com/OJOMB/fightpicker/internal/repo/fighters"
	usersrepo "github.com/OJOMB/fightpicker/internal/repo/users"
)

type repos struct {
	AuthRepo     *authrepo.Repo
	UsersRepo    *usersrepo.Repo
	FightersRepo *fightersrepo.Repo
	CardsRepo    *cardsrepo.Repo
}

func (a *App) newRepos(ctx context.Context, cfg *config.Config) error {
	authRepo, err := authrepo.New(a.DB.pool, a.DB.queries, a.Utils.IDTool, a.Utils.DateTimeTool, a.Logger)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create auth repo", "error", err)
		return err
	}

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
		return err
	}

	fightersRepo, err := fightersrepo.New(a.DB.pool, a.DB.queries, a.Clients.RedisClient, a.Utils.DateTimeTool, a.Logger)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create fighters repo", "error", err)
		return err
	}

	cardsRepo, err := cardsrepo.New(a.DB.pool, a.DB.queries, a.Logger)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create cards repo", "error", err)
		return err
	}

	a.Repos = repos{
		AuthRepo:     authRepo,
		UsersRepo:    usersRepo,
		FightersRepo: fightersRepo,
		CardsRepo:    cardsRepo,
	}

	return nil
}
