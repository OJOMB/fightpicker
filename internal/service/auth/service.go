package auth

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"

	serviceusers "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/jsonwebtokens"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// UserRepo defines the interface for user-related operations needed by the auth service.
type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (serviceusers.User, error)
	UpdateLastLoginAtByUserID(ctx context.Context, userID uuid.UUID) error
}

// AuthRepo defines the interface for authentication-related operations needed by the auth service.
type AuthRepo interface {
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, Permissions, error)
	StoreRefreshToken(ctx context.Context, userID, jti uuid.UUID, tokenHash, ipAddress, userAgent string, expiresAt time.Time) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	RotateRefreshTokens(ctx context.Context, oldTokenHash, newTokenHash string, newJTI, userID uuid.UUID, ipAddress, userAgent string, newExpiresAt time.Time) error
	RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error
}

// PasswordVerifier defines the interface for verifying passwords.
type PasswordVerifier interface {
	Verify(hashedPassword, password string) (bool, error)
}

// JWTGenerator defines the interface for generating and hashing JWT tokens.
type JWTGenerator interface {
	GenerateToken(jti, userID uuid.UUID, duration time.Duration, iss, aud string, customClaims map[string]any, secretKey []byte) (jsonwebtokens.Token, error)
	HashTokenString(tokenStr string) string
}

// Service provides authentication services.
type Service struct {
	userRepo        UserRepo
	authRepo        AuthRepo
	idGen           id.Generator
	pwordVerifier   PasswordVerifier
	jwtGen          JWTGenerator
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	secretKey       []byte
	tokenAudience   string
	tokenIssuer     string
	logger          logs.Logger
}

// New creates a new instance of the auth Service.
func New(userRepo UserRepo, authRepo AuthRepo, pwordVerifier PasswordVerifier, idGen id.Generator, jwtGen JWTGenerator, accessTokenTTL, refreshTokenTTL time.Duration, secretKey string, tokenAudience string, tokenIssuer string, logger logs.Logger) (*Service, error) {
	if userRepo == nil {
		return nil, ErrNilUserRepo
	}

	if authRepo == nil {
		return nil, ErrNilAuthRepo
	}

	if pwordVerifier == nil {
		return nil, ErrNilPasswordVerifier
	}

	if idGen == nil {
		return nil, ErrNilIDGenerator
	}

	if jwtGen == nil {
		return nil, ErrNilJWTGenerator
	}

	if accessTokenTTL <= 0 {
		return nil, ErrInvalidAccessTokenTTL
	}

	if refreshTokenTTL <= 0 {
		return nil, ErrInvalidRefreshTokenTTL
	}

	if secretKey == "" {
		return nil, ErrEmptySecretKey
	}

	return &Service{
		userRepo:        userRepo,
		authRepo:        authRepo,
		idGen:           idGen,
		pwordVerifier:   pwordVerifier,
		jwtGen:          jwtGen,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		secretKey:       []byte(secretKey),
		tokenAudience:   tokenAudience,
		tokenIssuer:     tokenIssuer,
		logger:          logger.With("component", "auth_service"),
	}, nil
}
