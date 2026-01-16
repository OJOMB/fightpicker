package contextual

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

var KeyRequestUserAgent Contextkey = "request_user_agent"
var KeyRequestHost Contextkey = "request_host"
var KeyRequestRemoteAddr Contextkey = "request_remote_addr"
var KeyRequestRouteName Contextkey = "request_route_name"

var KeyRequestParamUserID Contextkey = "user_id"
var KeyRequestParamFighterID Contextkey = "fighter_id"

func (c *ContextTool) WithRequestValues(ctx context.Context, r *http.Request) context.Context {
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

	// get current route name
	route := mux.CurrentRoute(r)
	if route != nil {
		routeName := route.GetName()
		ctx = context.WithValue(ctx, KeyRequestRouteName, routeName)
	}

	return ctx
}

func (c *ContextTool) GetRequestValues(ctx context.Context) map[string]string {
	values := make(map[string]string)

	if v, ok := ctx.Value(KeyRequestUserAgent).(string); ok {
		values[string(KeyRequestUserAgent)] = v
	}

	if v, ok := ctx.Value(KeyRequestHost).(string); ok {
		values[string(KeyRequestHost)] = v
	}

	if v, ok := ctx.Value(KeyRequestRemoteAddr).(string); ok {
		values[string(KeyRequestRemoteAddr)] = v
	}

	if v, ok := ctx.Value(KeyRequestRouteName).(string); ok {
		values[string(KeyRequestRouteName)] = v
	}

	if v, ok := ctx.Value(KeyRequestParamUserID).(string); ok {
		values[string(KeyRequestParamUserID)] = v
	}

	if v, ok := ctx.Value(KeyRequestParamFighterID).(string); ok {
		values[string(KeyRequestParamFighterID)] = v
	}

	return values
}
