package contextual

import (
	"context"
	"net/http"
	"slices"

	"github.com/gorilla/mux"

	"github.com/OJOMB/fightpicker/pkg/id"
)

// GetRequestID retrieves the request ID from the context.
func (c *ContextTool) GetRequestID(ctx context.Context) id.UUID7 {
	reqID, ok := ctx.Value(KeyRequestID).(string)
	if !ok {
		return id.UUID7Nil
	}

	uid, err := c.id.ParseString(reqID)
	if err != nil {
		return id.UUID7Nil
	}

	return uid
}

// GetReqSubjectFromContext retrieves the user ID (subject) of the user making the request from the context.
func (c *ContextTool) GetReqSubjectFromContext(ctx context.Context) id.UUID7 {
	subjectID, ok := ctx.Value(KeyRequestSubject).(string)
	if !ok {
		return id.UUID7Nil
	}

	uid, err := c.id.ParseString(subjectID)
	if err != nil {
		return id.UUID7Nil
	}

	return uid
}

// GetUserRolesFromContext retrieves the roles of the user making the request from the context.
func (c *ContextTool) GetUserRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(KeyUserRoles).([]string)
	return roles, ok
}

// ReqSubjectIsAnAdmin checks if the user making the request has the "admin" role.
func (c *ContextTool) ReqSubjectIsAnAdmin(ctx context.Context) bool {
	roles, ok := c.GetUserRolesFromContext(ctx)
	if !ok {
		return false
	}

	return slices.Contains(roles, "admin")
}

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
