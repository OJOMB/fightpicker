package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"

	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	"github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type JWTValidator interface {
	Parse(tokenStr string, secretKey []byte) (*auth.AuthClaims, *jwt.RegisteredClaims, error)
}

type AuthPermissionsChecker struct {
	logger       logs.Logger
	jwtValidator JWTValidator
	ignorePaths  map[string]struct{}
	secretKey    []byte
}

func NewAuthPermissionsChecker(secretKey []byte, jwtValidator JWTValidator, l logs.Logger) *AuthPermissionsChecker {
	return &AuthPermissionsChecker{
		secretKey:    secretKey,
		jwtValidator: jwtValidator,
		ignorePaths: map[string]struct{}{
			v1.EndpointNameV1AuthLogin:              {},
			v1.EndpointNameV1AuthRefresh:            {},
			v1.EndpointNameV1AuthLogout:             {},
			v1.EndpointNameV1UsersCreate:            {},
			v1.EndpointNameV1UsersEmailVerification: {},
		},
		logger: l.With("middleware", "AuthPermissionsChecker"),
	}
}

func (apc *AuthPermissionsChecker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		route := mux.CurrentRoute(r)
		name := route.GetName()

		permissionComponents := strings.Split(name, ".")
		if len(permissionComponents) != 4 {
			apc.logger.ErrorContext(ctx, "invalid endpoint name format", "endpoint_name", name)
			http.Error(w, "internal server error", http.StatusInternalServerError)
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
			req_ID, _ := r.Context().Value(contextual.KeyRequestID).(string)
			http.Error(w, fmt.Sprintf(`{
				"error":{
					"message":"missing authorization token",
					"code": "UNAUTHORIZED",
					"request_id": "%s"
				}
			}`, req_ID), http.StatusUnauthorized)
			return
		}

		// validate the token and permissions
		customClaims, registeredClaims, err := apc.jwtValidator.Parse(token, apc.secretKey)
		if err != nil {
			req_ID, _ := r.Context().Value(contextual.KeyRequestID).(string)
			http.Error(w, fmt.Sprintf(`{
				"error":{
					"message":"invalid authorization token",
					"code": "UNAUTHORIZED",
					"request_id": "%s"
				}
			}`, req_ID), http.StatusUnauthorized)
			return
		}

		// check the token is not expired
		if time.Now().After(registeredClaims.ExpiresAt.Time) {
			req_ID, _ := r.Context().Value(contextual.KeyRequestID).(string)
			http.Error(w, fmt.Sprintf(`{
				"error":{
					"message":"authorization token expired",
					"code": "UNAUTHORIZED",
					"request_id": "%s"
				}
			}`, req_ID), http.StatusUnauthorized)
			return
		}

		// check permissions
		if _, ok := customClaims.Perms[version][resource][operation][permissionName]; !ok {
			req_ID, _ := r.Context().Value(contextual.KeyRequestID).(string)
			http.Error(w, fmt.Sprintf(`{
				"error":{
					"message":"insufficient permissions",
					"code": "FORBIDDEN",
					"request_id": "%s"
				}
			}`, req_ID), http.StatusForbidden)
			return
		}

		// put the claims subject (user ID) into context
		subject, err := registeredClaims.GetSubject()
		if err != nil || subject == "" {
			req_ID, _ := r.Context().Value(contextual.KeyRequestID).(string)
			http.Error(w, fmt.Sprintf(`{
				"error":{
					"message":"invalid authorization token",
					"code": "UNAUTHORIZED",
					"request_id": "%s"
				}
			}`, req_ID), http.StatusUnauthorized)
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
