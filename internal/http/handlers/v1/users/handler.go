package users

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/OJOMB/fightpicker/internal/http/apiresponder"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/id"
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
	*v1.Handler
	service    Service
	pathPrefix string
}

func New(service Service, idTool id.UUID7Parser, ctxTool contextual.ContextProvider, logger logs.Logger) (*Handler, error) {
	if logger == nil {
		return nil, ErrLoggerIsNil
	}

	if idTool == nil {
		return nil, ErrIDToolIsNil
	}

	if ctxTool == nil {
		return nil, ErrContextToolIsNil
	}

	responder := apiresponder.NewJSONResponder(ctxTool, classifyError, logger.With("component", "handler_users_v1"))

	return &Handler{
		Handler:    v1.NewHandler(idTool, responder),
		service:    service,
		pathPrefix: pathPrefix,
	}, nil
}

func (h *Handler) RegisterRoutes(mux *mux.Router) {
	// POST /api/v1/users - create a new user
	h.AddRoute(mux, h.pathPrefix, http.MethodPost, v1.EndpointNameV1UsersCreate, h.createUser(h.service))

	// GET /api/v1/users/verify?token={token} - verification url sent out in email on user creation
	h.AddRoute(mux, fmt.Sprintf("%s/verify", h.pathPrefix), http.MethodGet, v1.EndpointNameV1UsersEmailVerification, h.verifyEmail(h.service))

	// GET /api/v1/users/{user_id} - get a user by ID
	h.AddRoute(mux, fmt.Sprintf("%s/{%s}", h.pathPrefix, v1.QueryParamUserID), http.MethodGet, v1.EndpointNameV1UsersGet, h.getUser(h.service))

	// GET /api/v1/users - list users with pagination
	h.AddRoute(mux, h.pathPrefix, http.MethodGet, v1.EndpointNameV1UsersList, h.listUsers(h.service))

	// PATCH /api/v1/users/{user_id} - update a user by ID
	h.AddRoute(mux, fmt.Sprintf("%s/{%s}", h.pathPrefix, v1.QueryParamUserID), http.MethodPatch, v1.EndpointNameV1UsersUpdate, h.updateUser(h.service))

	// DELETE /api/v1/users/{user_id} - delete a user by ID
	h.AddRoute(mux, fmt.Sprintf("%s/{%s}", h.pathPrefix, v1.QueryParamUserID), http.MethodDelete, v1.EndpointNameV1UsersDelete, h.deleteUser(h.service))

	// GET /api/v1/users/{user_id}/followers - list followers of a user by ID
	h.AddRoute(mux, fmt.Sprintf("%s/{%s}/followers", h.pathPrefix, v1.QueryParamUserID), http.MethodGet, v1.EndpointNameV1UsersListFollowers, h.listFollowers(h.service))

	// GET /api/v1/users/{user_id}/followees - list followees of a user by ID
	h.AddRoute(mux, fmt.Sprintf("%s/{%s}/followees", h.pathPrefix, v1.QueryParamUserID), http.MethodGet, v1.EndpointNameV1UsersListFollowees, h.listFollowees(h.service))

	// PUT /api/v1/users/{user_id}/follow - follow a user by ID
	h.AddRoute(mux, fmt.Sprintf("%s/{%s}/follow", h.pathPrefix, v1.QueryParamUserID), http.MethodPut, v1.EndpointNameV1UsersFollow, h.followUser(h.service))

	// DELETE /api/v1/users/{user_id}/unfollow - unfollow a user by ID
	h.AddRoute(mux, fmt.Sprintf("%s/{%s}/follow", h.pathPrefix, v1.QueryParamUserID), http.MethodDelete, v1.EndpointNameV1UsersUnfollow, h.unfollowUser(h.service))

	// POST /api/v1/users/{user_id}/profile-picture/upload-url - generate a presigned URL for user profile picture upload
	h.AddRoute(mux, fmt.Sprintf("%s/{%s}/profile-picture/upload-url", h.pathPrefix, v1.QueryParamUserID), http.MethodPost, v1.EndpointNameV1UsersGeneratePresignedURL, h.generatePresignedURL(h.service))
}
