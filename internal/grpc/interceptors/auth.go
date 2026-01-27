package interceptors

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/contextual"
)

type JWTValidator interface {
	Parse(tokenStr string, secretKey []byte) (*auth.AuthClaims, *jwt.RegisteredClaims, error)
}

type UnaryAuthInterceptor struct {
	secretKey     []byte
	jwt           JWTValidator
	ctxTool       contextual.ContextProvider
	ignoreMethods map[string]struct{}
}

func NewUnaryAuthInterceptor(secretKey []byte, jwtValidator JWTValidator, ctxTool contextual.ContextProvider) (*UnaryAuthInterceptor, error) {
	if len(secretKey) == 0 {
		return nil, ErrSecretKeyIsNilOrEmpty
	}

	if jwtValidator == nil {
		return nil, ErrJWTValidatorIsNil
	}

	if ctxTool == nil {
		return nil, ErrContextToolIsNil
	}

	return &UnaryAuthInterceptor{
		secretKey: secretKey,
		jwt:       jwtValidator,
		ctxTool:   ctxTool,
		ignoreMethods: map[string]struct{}{
			"/v1.Auth/Login":        {},
			"/v1.Auth/Refresh":      {},
			"/v1.Auth/Logout":       {},
			"/v1.Users/Create":      {},
			"/v1.Users/VerifyEmail": {},
		},
	}, nil
}

func (i UnaryAuthInterceptor) Intercept() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		fullMethod := info.FullMethod
		// Bypass auth for certain methods
		if _, ok := i.ignoreMethods[fullMethod]; ok {
			return handler(ctx, req)
		}

		// 1. Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is missing")
		}

		// 2. Look for the "authorization" key
		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is missing")
		}

		// 3. Validate the token (using your existing pkg/jsonwebtokens)
		tokenString := strings.TrimPrefix(values[0], "Bearer ")
		claims, registeredClaims, err := i.jwt.Parse(tokenString, i.secretKey)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		if time.Now().After(registeredClaims.ExpiresAt.Time) {
			return nil, status.Errorf(codes.Unauthenticated, "token has expired")
		}

		subject, err := registeredClaims.GetSubject()
		if err != nil || subject == "" {
			return nil, status.Errorf(codes.Unauthenticated, "error getting token subject: %v", err)
		}

		// 4. Add claims to context
		claimsCtx := context.WithValue(ctx, contextual.KeyRequestSubject, subject)
		claimsCtx = context.WithValue(claimsCtx, contextual.KeyUserRoles, claims.Roles)

		return handler(claimsCtx, req)
	}
}
