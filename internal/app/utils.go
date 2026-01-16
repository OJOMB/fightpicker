package app

import (
	"github.com/pkg/errors"

	"github.com/OJOMB/fightpicker/internal/config"
	serviceauth "github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/auth"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/jsonwebtokens"
	"github.com/OJOMB/fightpicker/pkg/media"
)

type utils struct {
	IDTool         id.UUID7GeneratorParser
	DateTimeTool   datetimes.Now
	AuthTool       *auth.AuthTool
	ImageProcessor *media.ImageProcessor
	JWTTool        *jsonwebtokens.JWTTool[serviceauth.AuthClaims]
	ContextTool    contextual.ContextProvider
}

func (a *App) newUtils(cfg *config.Config) error {
	idTool := id.NewUUIDV7Tool()
	dateTimeTool := datetimes.NewUTCNow()

	if cfg.Auth.HashingCost <= 0 {
		return errors.Wrap(ErrInvalidConfig, "hashing_cost")
	}
	authTool := auth.NewAuthTool(cfg.Auth.HashingCost)

	imageProcessor := media.NewImageProcessor()

	jwtTool := jsonwebtokens.NewJWTTool[serviceauth.AuthClaims]()

	ctxTool := contextual.NewContextTool(idTool)

	a.Utils = &utils{
		IDTool:         idTool,
		DateTimeTool:   dateTimeTool,
		AuthTool:       authTool,
		ImageProcessor: imageProcessor,
		JWTTool:        jwtTool,
		ContextTool:    ctxTool,
	}

	return nil
}
