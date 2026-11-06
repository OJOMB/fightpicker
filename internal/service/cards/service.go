package cards

import (
	"context"

	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type Repo interface {
	CreateCard(ctx context.Context, user Card) error
}

type Service struct {
	repo         Repo
	idTool       id.UUID7GeneratorParser
	dateTimeTool datetimes.Now
	ctxTool      contextual.ContextProvider
	logger       logs.Logger
}

func New(repo Repo, datetimeTool datetimes.Now, ctxTool contextual.ContextProvider, idTool id.UUID7GeneratorParser, logger logs.Logger) (*Service, error) {
	if repo == nil {
		return nil, errNilRepo
	}

	if datetimeTool == nil {
		return nil, errNilDatetimeTool
	}

	if ctxTool == nil {
		return nil, errNilContextTool
	}

	if idTool == nil {
		return nil, errNilIDTool
	}

	if logger == nil {
		return nil, errNilLogger
	}

	return &Service{repo: repo, dateTimeTool: datetimeTool, ctxTool: ctxTool, idTool: idTool, logger: logger}, nil
}
