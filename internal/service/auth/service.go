package auth

import (
	"context"
	"time"

	serviceusers "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/jsonwebtokens"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// UserRepo defines the interface for user-related operations needed by the auth service.
type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (serviceusers.User, error)
	UpdateLastLoginAtByUserID(ctx context.Context, userID id.UUID7) error
}

// AuthRepo defines the interface for authentication-related operations needed by the auth service.
type AuthRepo interface {
	GetUserPermissions(ctx context.Context, userID id.UUID7) ([]string, Permissions, error)
	StoreRefreshToken(ctx context.Context, userID, jti id.UUID7, tokenHash, ipAddress, userAgent string, expiresAt time.Time) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	RotateRefreshTokens(ctx context.Context, oldTokenHash, newTokenHash string, newJTI, userID id.UUID7, ipAddress, userAgent string, newExpiresAt time.Time) error
	RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error
}

// PasswordVerifier defines the interface for verifying passwords.
type PasswordVerifier interface {
	Verify(hashedPassword, password string) (bool, error)
}

// JWTGenerator defines the interface for generating and hashing JWT tokens.
type JWTGenerator interface {
	GenerateToken(jti, userID id.UUID7, duration time.Duration, iss, aud string, customClaims map[string]any, secretKey []byte) (jsonwebtokens.Token, error)
	HashTokenString(tokenStr string) string
}

// Service provides authentication services.
type Service struct {
	userRepo         UserRepo
	authRepo         AuthRepo
	id               id.UUID7GeneratorParser
	passwordVerifier PasswordVerifier
	jwtGen           JWTGenerator
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	secretKey        []byte
	tokenAudience    string
	tokenIssuer      string
	logger           logs.Logger
}

// New creates a new instance of the auth Service.
func New(userRepo UserRepo, authRepo AuthRepo, passwordVerifier PasswordVerifier, idTool id.UUID7GeneratorParser, jwtGen JWTGenerator, accessTokenTTL, refreshTokenTTL time.Duration, secretKey string, tokenAudience string, tokenIssuer string, logger logs.Logger) (*Service, error) {
	if userRepo == nil {
		return nil, errNilUserRepo
	}

	if authRepo == nil {
		return nil, errNilAuthRepo
	}

	if passwordVerifier == nil {
		return nil, errNilPasswordVerifier
	}

	if idTool == nil {
		return nil, errNilIDGenerator
	}

	if jwtGen == nil {
		return nil, errNilJWTGenerator
	}

	if accessTokenTTL <= 0 {
		return nil, errInvalidAccessTokenTTL
	}

	if refreshTokenTTL <= 0 {
		return nil, errInvalidRefreshTokenTTL
	}

	if secretKey == "" {
		return nil, errEmptySecretKey
	}

	return &Service{
		userRepo:         userRepo,
		authRepo:         authRepo,
		id:               idTool,
		passwordVerifier: passwordVerifier,
		jwtGen:           jwtGen,
		accessTokenTTL:   accessTokenTTL,
		refreshTokenTTL:  refreshTokenTTL,
		secretKey:        []byte(secretKey),
		tokenAudience:    tokenAudience,
		tokenIssuer:      tokenIssuer,
		logger:           logger.With("component", "auth_service"),
	}, nil
}
