package fighters

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type Repo interface {
	GetFighterByID(ctx context.Context, id id.UUID7) (Fighter, error)
}

type Service struct {
	repo   Repo
	logger logs.Logger
}

func New(repo Repo, idGen id.UUID7GeneratorParser, now datetimes.Now, logger logs.Logger) (*Service, error) {
	if repo == nil {
		return nil, ErrNilRepo
	}

	if logger == nil {
		return nil, ErrNilLogger
	}

	if idGen == nil {
		return nil, ErrNilIDGenerator
	}

	if now == nil {
		return nil, ErrNilNowFunc
	}

	return &Service{
		repo:   repo,
		logger: logger.With("component", "fighters_service"),
	}, nil
}
