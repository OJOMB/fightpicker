package app

import (
	"context"

	"github.com/OJOMB/fightpicker/internal/consumers/email"
	"github.com/OJOMB/fightpicker/internal/consumers/media"
)

// RunWorkers starts all background worker consumers.
func (a *App) RunWorkers(ctx context.Context) error {
	if !a.initialized {
		return ErrAppNotInitialized
	}

	// create consumers
	imageConsumer, err := media.NewUserProfilePictureConsumer(
		a.Clients.KafkaMediaProfilePicClient,
		a.Services.UsersService,
		a.Logger,
	)
	if err != nil {
		a.Logger.FatalContext(ctx, "failed to create user profile picture consumer", "error", err)
	}

	emailConsumer, err := email.NewUserCreationEmailConsumer(
		a.Clients.KafkaPostUserCreateClient,
		a.Services.UsersService,
		a.Logger,
	)
	if err != nil {
		a.Logger.FatalContext(ctx, "failed to create user creation email consumer", "error", err)
	}

	// kick off the consumers in separate goroutines
	go imageConsumer.Run(ctx)
	go emailConsumer.Run(ctx)

	return nil
}
