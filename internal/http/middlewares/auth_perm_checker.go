package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"

	"github.com/OJOMB/fightpicker/internal/http/apiresponder"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type JWTValidator interface {
	Parse(tokenStr string, secretKey []byte) (*auth.AuthClaims, *jwt.RegisteredClaims, error)
}

type AuthPermissionsChecker struct {
	apiresponder.Responder
	jwtValidator JWTValidator
	ctxTool      contextual.ContextTool
	ignorePaths  map[string]struct{}
	secretKey    []byte
	logger       logs.Logger
}

func NewAuthPermissionsChecker(secretKey []byte, jwtValidator JWTValidator, ctxTool contextual.ContextProvider, l logs.Logger) (*AuthPermissionsChecker, error) {
	if secretKey == nil || len(secretKey) == 0 {
		return nil, ErrSecretKeyIsNilOrEmpty
	}

	if jwtValidator == nil {
		return nil, ErrJWTValidatorIsNil
	}

	if ctxTool == nil {
		return nil, ErrContextToolIsNil
	}

	if l == nil {
		return nil, ErrLoggerIsNil
	}

	return &AuthPermissionsChecker{
		Responder:    apiresponder.NewJSONResponder(ctxTool, classifyError, l.With("component", "middleware_auth_perm_checker")),
		secretKey:    secretKey,
		jwtValidator: jwtValidator,
		ignorePaths: map[string]struct{}{
			v1.EndpointNameV1AuthLogin:              {},
			v1.EndpointNameV1AuthRefresh:            {},
			v1.EndpointNameV1AuthLogout:             {},
			v1.EndpointNameV1UsersCreate:            {},
			v1.EndpointNameV1UsersEmailVerification: {},
		},
	}, nil
}

func (apc *AuthPermissionsChecker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		route := mux.CurrentRoute(r)
		name := route.GetName()

		permissionComponents := strings.Split(name, ".")
		if len(permissionComponents) != 4 {
			apc.WriteError(ctx, w, ErrInvalidRouteName)
			return
		}

		version := permissionComponents[0]
		resource := permissionComponents[1]
		operation := permissionComponents[2]
		permissionName := permissionComponents[3]

		// skip checking for ignored paths that don't need auth
		if _, ok := apc.ignorePaths[name]; ok || version == "static" {
			next.ServeHTTP(w, r)
			return
		}

		// get the token from header
		bearerToken := r.Header.Get("Authorization")
		token := strings.TrimPrefix(bearerToken, "Bearer ")
		if token == "" {
			apc.WriteError(ctx, w, ErrMissingToken)
			return
		}

		// validate the token and permissions
		customClaims, registeredClaims, err := apc.jwtValidator.Parse(token, apc.secretKey)
		if err != nil {
			apc.WriteError(ctx, w, ErrInvalidToken)
			return
		}

		// check the token is not expired
		if time.Now().After(registeredClaims.ExpiresAt.Time) {
			apc.WriteError(ctx, w, ErrExpiredToken)
			return
		}

		// check permissions
		if _, ok := customClaims.Perms[version][resource][operation][permissionName]; !ok {
			apc.WriteError(ctx, w, ErrInsufficientPermissions)
			return
		}

		// put the claims subject (user ID) into context
		subject, err := registeredClaims.GetSubject()
		if err != nil || subject == "" {
			apc.WriteError(ctx, w, ErrInvalidToken)
			return
		}

		ctx = context.WithValue(ctx, contextual.KeyRequestSubject, subject)
		r = r.WithContext(ctx)

		// add user roles to context
		var roles = make([]string, len(customClaims.Roles))
		copy(roles, customClaims.Roles)
		ctx = context.WithValue(ctx, contextual.KeyUserRoles, roles)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
