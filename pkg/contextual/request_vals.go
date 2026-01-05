package contextual

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

var KeyRequestUserAgent Contextkey = "request_user_agent"
var KeyRequestHost Contextkey = "request_host"
var KeyRequestRemoteAddr Contextkey = "request_remote_addr"

func WithRequestValues(ctx context.Context, r *http.Request) context.Context {
	// get User-Agent header
	userAgent := r.Header.Get("User-Agent")
	ctx = context.WithValue(ctx, KeyRequestUserAgent, userAgent)

	// get server Host
	ctx = context.WithValue(ctx, KeyRequestHost, r.Host)
	// get remote address
	ctx = context.WithValue(ctx, KeyRequestRemoteAddr, r.RemoteAddr)

	// get path variables
	vars := mux.Vars(r)
	for k, v := range vars {
		ctx = context.WithValue(ctx, Contextkey(k), v)
	}

	return ctx
}
