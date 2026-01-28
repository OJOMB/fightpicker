package auth

import "fmt"

var (
	///////////////////////////////
	// initialization errors     //
	///////////////////////////////

	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = fmt.Errorf("missing required parameter")

	// ErrNilUserService is returned when the provided user service is nil.
	ErrNilUserRepo = fmt.Errorf("user repository is nil")
	// ErrNilAuthService is returned when the provided auth service is nil.
	ErrNilAuthRepo = fmt.Errorf("auth repository is nil")

	// ErrNilPasswordHasher is returned when the password hasher is nil.
	ErrNilPasswordHasher = fmt.Errorf("password hasher is nil")
	// ErrNilTokenStore is returned when the token store is nil.
	ErrNilTokenStore = fmt.Errorf("token store is nil")
	// ErrNilPasswordVerifier is returned when the password verifier is nil.
	ErrNilPasswordVerifier = fmt.Errorf("password verifier is nil")
	// ErrNilTimeProvider is returned when the time provider is nil.
	ErrNilTimeProvider = fmt.Errorf("time provider is nil")
	// ErrNilIDGenerator is returned when the ID generator is nil.
	ErrNilIDGenerator = fmt.Errorf("ID generator is nil")
	// ErrNilJWTGenerator is returned when the JWT generator is nil.
	ErrNilJWTGenerator = fmt.Errorf("JWT generator is nil")

	// ErrInvalidAccessTokenTTL is returned when the access token TTL is invalid.
	ErrInvalidAccessTokenTTL = fmt.Errorf("access token TTL must be greater than zero")
	// ErrInvalidRefreshTokenTTL is returned when the refresh token TTL is invalid.
	ErrInvalidRefreshTokenTTL = fmt.Errorf("refresh token TTL must be greater than zero")

	// ErrEmptySecretKey is returned when the secret key is empty.
	ErrEmptySecretKey = fmt.Errorf("secret key cannot be empty")

	// ErrInvalidCredentials is returned when the provided credentials are invalid.
	ErrInvalidCredentials = fmt.Errorf("invalid credentials provided")

	///////////////////////////////
	// runtime errors            //
	///////////////////////////////

	// ErrUserNotFound is returned when a user is not found in the repository.
	ErrUserNotFound = fmt.Errorf("user not found")

	// ErrRefreshTokenRevoked is returned when a refresh token has been revoked.
	ErrRefreshTokenRevoked = fmt.Errorf("refresh token has been revoked")

	// ErrRefreshTokenExpired is returned when a refresh token has expired.
	ErrRefreshTokenExpired = fmt.Errorf("refresh token has expired")

	// ErrRefreshTokenReused is returned when a refresh token has been reused.
	ErrRefreshTokenReused = fmt.Errorf("refresh token has been reused")
)
