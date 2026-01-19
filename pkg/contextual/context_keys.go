package contextual

// Contextkey is a type for context keys used in this package.
type Contextkey string

var (
	// KeyRequestID is the context key for the unique request ID.
	KeyRequestID Contextkey = "request_id"

	// KeyRequestUserAgent is the context key for the User-Agent header of the request.
	KeyRequestUserAgent Contextkey = "request_user_agent"

	// KeyRequestHost is the context key for the host of the request.
	KeyRequestHost Contextkey = "request_host"

	// KeyRequestRemoteAddr is the context key for the remote address of the request.
	KeyRequestRemoteAddr Contextkey = "request_remote_addr"

	// KeyRequestRouteName is the context key for the current route name.
	KeyRequestRouteName Contextkey = "request_route_name"

	// KeyRequestParamUserID is the context key for the user ID path parameter.
	KeyRequestParamUserID Contextkey = "user_id"

	// KeyRequestParamFighterID is the context key for the fighter ID path parameter.
	KeyRequestParamFighterID Contextkey = "fighter_id"

	// KeyRequestSubject is the context key for the user ID (subject) of the user making the request.
	KeyRequestSubject Contextkey = "request_subject"

	// KeyUserRoles is the context key for the roles of the user making the request.
	KeyUserRoles Contextkey = "user_roles"
)
