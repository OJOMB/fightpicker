package auth

import "fmt"

var (
	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = fmt.Errorf("missing required parameter")

	ErrNilUserRepo = fmt.Errorf("user repository is nil")
	ErrNilAuthRepo = fmt.Errorf("auth repository is nil")

	ErrNilPasswordVerifier = fmt.Errorf("password verifier is nil")
	ErrNilIDGenerator      = fmt.Errorf("ID generator is nil")
	ErrNilJWTGenerator     = fmt.Errorf("JWT generator is nil")

	ErrInvalidAccessTokenTTL  = fmt.Errorf("access token TTL must be greater than zero")
	ErrInvalidRefreshTokenTTL = fmt.Errorf("refresh token TTL must be greater than zero")

	ErrEmptySecretKey = fmt.Errorf("secret key cannot be empty")

	// ErrInvalidCredentials is returned when the provided credentials are invalid.
	ErrInvalidCredentials = fmt.Errorf("invalid credentials provided")

	// ErrUserNotFound is returned when a user is not found in the repository.
	ErrUserNotFound = fmt.Errorf("user not found")

	// ErrRefreshTokenRevoked is returned when a refresh token has been revoked.
	ErrRefreshTokenRevoked = fmt.Errorf("refresh token has been revoked")

	// ErrRefreshTokenExpired is returned when a refresh token has expired.
	ErrRefreshTokenExpired = fmt.Errorf("refresh token has expired")

	// ErrRefreshTokenReused is returned when a refresh token has been reused.
	ErrRefreshTokenReused = fmt.Errorf("refresh token has been reused")
)
