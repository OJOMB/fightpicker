package fighters

import (
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type Repo interface {
	FighterByIDGetter
	FightersIngestor
}

type Service struct {
	repo         Repo
	id           id.UUID7GeneratorParser
	datetimeTool datetimes.Now
	ctxTool      contextual.ContextProvider
	logger       logs.Logger
}

func New(repo Repo, idGen id.UUID7GeneratorParser, now datetimes.Now, ctxProvider contextual.ContextProvider, logger logs.Logger) (*Service, error) {
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
		id:           idGen,
		datetimeTool: now,
		ctxTool:      ctxProvider,
		repo:         repo,
		logger:       logger.With("component", "fighters_service"),
	}, nil
}
