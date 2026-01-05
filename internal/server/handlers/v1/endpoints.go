package v1

const (
	// v1 Auth
	EndpointNameV1AuthLogin   = "v1.auth.post.login"
	EndpointNameV1AuthRefresh = "v1.auth.post.refresh"
	EndpointNameV1AuthLogout  = "v1.auth.post.logout"

	// v1 Users
	EndpointNameV1UsersCreate               = "v1.users.post.create"
	EndpointNameV1UsersEmailVerification    = "v1.users.get.verify_email"
	EndpointNameV1UsersGet                  = "v1.users.get.get"
	EndpointNameV1UsersList                 = "v1.users.get.list"
	EndpointNameV1UsersListFollowers        = "v1.users.get.list_followers"
	EndpointNameV1UsersListFollowees        = "v1.users.get.list_followees"
	EndpointNameV1UsersUpdate               = "v1.users.patch.update"
	EndpointNameV1UsersDelete               = "v1.users.delete.delete"
	EndpointNameV1UsersFollow               = "v1.users.put.follow"
	EndpointNameV1UsersUnfollow             = "v1.users.delete.unfollow"
	EndpointNameV1UsersGeneratePresignedURL = "v1.users.post.generate_presigned_url"

	// v1 Fighters
	EndpointNameV1FightersCreate = "v1.fighters.post.create"
	EndpointNameV1FightersGet    = "v1.fighters.get.get"
	EndpointNameV1FightersList   = "v1.fighters.get.list"
	EndpointNameV1FightersUpdate = "v1.fighters.patch.update"
	EndpointNameV1FightersDelete = "v1.fighters.delete.delete"
)
