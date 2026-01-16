package app

import (
	"context"
	"fmt"
	"time"

	"github.com/OJOMB/fightpicker/internal/config"
	authservice "github.com/OJOMB/fightpicker/internal/service/auth"
	fightersservice "github.com/OJOMB/fightpicker/internal/service/fighters"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
)

type services struct {
	AuthService     *authservice.Service
	UsersService    *usersservice.Service
	FightersService *fightersservice.Service
}

func (a *App) newServices(ctx context.Context, cfg *config.Config) error {
	// auth Service + Handler
	accessTokenTTL := time.Hour * time.Duration(cfg.Auth.AccessTTLHours)
	refreshTokenTTL := time.Hour * time.Duration(cfg.Auth.RefreshTTLHours)

	authService, err := authservice.New(
		a.Repos.UsersRepo,
		a.Repos.AuthRepo,
		a.Utils.AuthTool,
		a.Utils.IDTool,
		a.Utils.JWTTool,
		accessTokenTTL,
		refreshTokenTTL,
		cfg.Auth.PrivateKey,
		cfg.Auth.TokenAudience,
		cfg.Auth.TokenIssuer,
		a.Logger,
	)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create auth service", "error", err)
		return err
	}

	usersService, err := usersservice.New(
		a.Repos.UsersRepo,
		a.Utils.IDTool,
		a.Utils.DateTimeTool,
		a.Utils.ContextTool,
		a.Utils.AuthTool,
		a.Utils.ImageProcessor,
		fmt.Sprintf("http://%s:%d", cfg.HTTP.Domain, cfg.HTTP.Port),
		cfg.Email.AddressNoReply,
		a.Logger,
	)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create user service", "error", err)
		return err
	}

	fightersService, err := fightersservice.New(a.Repos.FightersRepo, a.Utils.IDTool, a.Utils.DateTimeTool, a.Logger)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed to create fighters service", "error", err)
		return err
	}

	a.Services = services{
		AuthService:     authService,
		UsersService:    usersService,
		FightersService: fightersService,
	}

	return nil
}
