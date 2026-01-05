package contextual

import (
	"context"
	"slices"

	"github.com/gofrs/uuid"
)

// Contextkey is a type for context keys used in this package.
type Contextkey string

// KeyRequestID is the context key for the current request ID.
var KeyRequestID Contextkey = "request_id"

// KeyRequestSubject is the context key for the user ID (subject) of the user making the request.
var KeyRequestSubject Contextkey = "request_subject"

// KeyUserRoles is the context key for the roles of the user making the request.
var KeyUserRoles Contextkey = "user_roles"

// GetReqSubjectFromContext retrieves the user ID (subject) of the user making the request from the context.
func GetReqSubjectFromContext(ctx context.Context) uuid.UUID {
	subjectID, ok := ctx.Value(KeyRequestSubject).(string)
	if !ok {
		return uuid.UUID{}
	}

	id, err := uuid.FromString(subjectID)
	if err != nil {
		return uuid.UUID{}
	}

	return id
}

// GetUserRolesFromContext retrieves the roles of the user making the request from the context.
func GetUserRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(KeyUserRoles).([]string)
	return roles, ok
}

// ReqSubjectIsAnAdmin checks if the user making the request has the "admin" role.
func ReqSubjectIsAnAdmin(ctx context.Context) bool {
	roles, ok := GetUserRolesFromContext(ctx)
	if !ok {
		return false
	}

	return slices.Contains(roles, "admin")
}
