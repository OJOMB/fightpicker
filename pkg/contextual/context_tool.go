package contextual

import (
	"context"
	"net/http"

	"github.com/OJOMB/fightpicker/pkg/id"
)

type ContextProvider interface {
	ContextServiceProvider
	ContextRequestProvider
}

type ContextServiceProvider interface {
	GetReqSubjectFromContext(ctx context.Context) id.UUID7
	GetUserRolesFromContext(ctx context.Context) ([]string, bool)
	ReqSubjectIsAnAdmin(ctx context.Context) bool
}

type ContextRequestProvider interface {
	WithRequestValues(ctx context.Context, r *http.Request) context.Context
	GetRequestValues(ctx context.Context) map[string]string
}

type ContextTool struct {
	id id.UUID7GeneratorParser
}

func NewContextTool(idTool id.UUID7GeneratorParser) *ContextTool {
	return &ContextTool{
		id: idTool,
	}
}
