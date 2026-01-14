package users

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

const pathPrefix = "/api/v1/users"

type Service interface {
	UserCreator
	UserGetter
	UserDeleter
	UserUpdater
	UserSearcher
	UserFollowerLister
	UserFolloweeLister
	UserFollower
	UserUnfollower
	PresignedPutURLGenerator
	EmailVerifier
}

type Handler struct {
	v1.Handler
	service    Service
	pathPrefix string
	logger     logs.Logger
}

func New(service Service, logger logs.Logger) *Handler {
	return &Handler{
		Handler: v1.Handler{
			Logger: logger.With("component", "http_handler_v1_users"),
		},
		service:    service,
		pathPrefix: pathPrefix,
	}
}

func (h *Handler) RegisterRoutes(mux *mux.Router) {
	// POST /api/v1/users - create a new user
	mux.Handle(
		h.pathPrefix,
		otelhttp.NewHandler(
			h.ToHandler(
				h.createUser(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersCreate,
		),
	).Name(v1.EndpointNameV1UsersCreate).
		Methods(http.MethodPost)

	// GET /api/v1/users/verify?token={token} - verification url sent out in email on user creation
	mux.Handle(
		fmt.Sprintf("%s/verify", h.pathPrefix),
		otelhttp.NewHandler(
			h.ToHandler(
				h.verifyEmail(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersEmailVerification,
		),
	).Name(v1.EndpointNameV1UsersEmailVerification).
		Methods(http.MethodGet)

	// GET /api/v1/users/{user_id} - get a user by ID
	mux.Handle(
		fmt.Sprintf("%s/{%s}", h.pathPrefix, v1.QueryParamUserID),
		otelhttp.NewHandler(
			h.ToHandler(
				h.getUser(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersGet,
		),
	).Name(v1.EndpointNameV1UsersGet).
		Methods(http.MethodGet)

	// GET /api/v1/users - list users with pagination
	mux.Handle(
		h.pathPrefix,
		otelhttp.NewHandler(
			h.ToHandler(
				h.listUsers(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersList,
		),
	).Name(v1.EndpointNameV1UsersList).
		Methods(http.MethodGet)

	// PATCH /api/v1/users/{user_id} - update a user by ID
	mux.Handle(
		fmt.Sprintf("%s/{%s}", h.pathPrefix, v1.QueryParamUserID),
		otelhttp.NewHandler(
			h.ToHandler(
				h.updateUser(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersUpdate,
		),
	).Name(v1.EndpointNameV1UsersUpdate).
		Methods(http.MethodPatch)

	// DELETE /api/v1/users/{user_id} - delete a user by ID
	mux.Handle(
		fmt.Sprintf("%s/{%s}", h.pathPrefix, v1.QueryParamUserID),
		otelhttp.NewHandler(
			h.ToHandler(
				h.deleteUser(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersDelete,
		),
	).Name(v1.EndpointNameV1UsersDelete).
		Methods(http.MethodDelete)

	// GET /api/v1/users/{user_id}/followers - list followers of a user by ID
	mux.Handle(
		fmt.Sprintf("%s/{%s}/followers", h.pathPrefix, v1.QueryParamUserID),
		otelhttp.NewHandler(
			h.ToHandler(
				h.listFollowers(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersListFollowers,
		),
	).Name(v1.EndpointNameV1UsersListFollowers).
		Methods(http.MethodGet)

	// GET /api/v1/users/{user_id}/followees - list followees of a user by ID
	mux.Handle(
		fmt.Sprintf("%s/{%s}/followees", h.pathPrefix, v1.QueryParamUserID),
		otelhttp.NewHandler(
			h.ToHandler(
				h.listFollowees(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersListFollowees,
		),
	).Name(v1.EndpointNameV1UsersListFollowees).
		Methods(http.MethodGet)

	// PUT /api/v1/users/{user_id}/follow - follow a user by ID
	mux.Handle(
		fmt.Sprintf("%s/{%s}/follow", h.pathPrefix, v1.QueryParamUserID),
		otelhttp.NewHandler(
			h.ToHandler(
				h.followUser(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersFollow,
		),
	).Name(v1.EndpointNameV1UsersFollow).
		Methods(http.MethodPut)

	// DELETE /api/v1/users/{user_id}/unfollow - unfollow a user by ID
	mux.Handle(
		fmt.Sprintf("%s/{%s}/follow", h.pathPrefix, v1.QueryParamUserID),
		otelhttp.NewHandler(
			h.ToHandler(
				h.unfollowUser(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersUnfollow,
		),
	).Name(v1.EndpointNameV1UsersUnfollow).
		Methods(http.MethodDelete)

	// POST /api/v1/users/{user_id}/profile-picture/upload-url - generate a presigned URL for user profile picture upload
	mux.Handle(
		fmt.Sprintf("%s/{%s}/profile-picture/upload-url", h.pathPrefix, v1.QueryParamUserID),
		otelhttp.NewHandler(
			h.ToHandler(
				h.generatePresignedURL(h.service),
				classifyError,
			),
			v1.EndpointNameV1UsersGeneratePresignedURL,
		),
	).Name(v1.EndpointNameV1UsersGeneratePresignedURL).
		Methods(http.MethodPost)
}
