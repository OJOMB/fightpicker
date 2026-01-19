package apiresponder

import (
	"context"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/apierr"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type ResponderBuilder interface {
	Build(errorClassifier apierr.APIErrClassifier, logger logs.Logger) (Responder, error)
}

type Responder interface {
	WriteError(ctx context.Context, w http.ResponseWriter, err error)
	Write(ctx context.Context, w http.ResponseWriter, status int, v any)
}
