package v1

var (
	// Query parameter names
	QueryParamPageSize               = "page_size"
	QueryParamLastSeenID             = "last_seen_id"
	QueryParamEmail                  = "email"
	QueryParamUserID                 = "user_id"
	QueryParamFighterID              = "fighter_id"
	QueryParamUsername               = "username"
	QueryParamEmailVerificationToken = "token"

	// Cookie names
	CookieKeyRefreshToken = "refresh_token"

	// Default values
	DefaultPageSize = 20
	MaxPageSize     = 100
)
