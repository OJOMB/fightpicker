package auth

import "fmt"

// Initialization errors (lowercase) indicate programming errors during
// service construction and are not exposed as they cannot be meaningfully handled at runtime.

// private initialization errors
var (
	// errNilUserRepo is returned when the provided user repository is nil.
	errNilUserRepo = fmt.Errorf("user repository is nil")
	// errNilAuthRepo is returned when the provided auth repository is nil.
	errNilAuthRepo = fmt.Errorf("auth repository is nil")

	// errNilPasswordVerifier is returned when the password verifier is nil.
	errNilPasswordVerifier = fmt.Errorf("password verifier is nil")
	// errNilTimeProvider is returned when the time provider is nil.
	errNilTimeProvider = fmt.Errorf("time provider is nil")
	// errNilIDGenerator is returned when the ID generator is nil.
	errNilIDGenerator = fmt.Errorf("ID generator is nil")
	// errNilJWTGenerator is returned when the JWT generator is nil.
	errNilJWTGenerator = fmt.Errorf("JWT generator is nil")
	// errEmptySecretKey is returned when the secret key is empty.
	errEmptySecretKey = fmt.Errorf("secret key cannot be empty")
	// errInvalidRefreshTokenTTL is returned when the refresh token TTL is invalid.
	errInvalidRefreshTokenTTL = fmt.Errorf("refresh token TTL must be greater than zero")
	// errInvalidAccessTokenTTL is returned when the access token TTL is invalid.
	errInvalidAccessTokenTTL = fmt.Errorf("access token TTL must be greater than zero")
)

// public runtime errors - these are the errors that the running transport layer might encounter and wish to handle
var (
	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = fmt.Errorf("missing required parameter")

	// ErrUserNotFound is returned when a user is not found in the repository.
	ErrUserNotFound = fmt.Errorf("user not found")

	// ErrRefreshTokenRevoked is returned when a refresh token has been revoked.
	ErrRefreshTokenRevoked = fmt.Errorf("refresh token has been revoked")

	// ErrRefreshTokenExpired is returned when a refresh token has expired.
	ErrRefreshTokenExpired = fmt.Errorf("refresh token has expired")

	// ErrRefreshTokenReused is returned when a refresh token has been reused.
	ErrRefreshTokenReused = fmt.Errorf("refresh token has been reused")
	// ErrInvalidCredentials is returned when the provided credentials are invalid.
	ErrInvalidCredentials = fmt.Errorf("invalid credentials provided")
)
