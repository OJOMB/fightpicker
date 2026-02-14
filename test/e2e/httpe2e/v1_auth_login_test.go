package httpe2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
)

func TestV1Login(t *testing.T) {
	// create a test user to login with
	email := newRandomEmail()
	password := newRandomString(12)
	user := createTestUser(t, email, password)
	defer cleanupTestUser(t, user.Id)

	type testCase struct {
		name               string
		requestBody        dtos.LoginRequest
		expectedStatusCode int
		expectedError      dtos.ErrorEnvelope
	}

	testCases := []testCase{
		{
			name: "successful login",
			requestBody: dtos.LoginRequest{
				Email:    openapi_types.Email(email),
				Password: password,
			},
			expectedStatusCode: 200,
		},
		{
			name: "missing password",
			requestBody: dtos.LoginRequest{
				Email: openapi_types.Email(email),
			},
			expectedStatusCode: 400,
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:    "MISSING_REQUIRED_PARAMETER",
					Message: "password: missing required parameter",
				},
			},
		},
		{
			name: "invalid email format",
			requestBody: dtos.LoginRequest{
				Email:    openapi_types.Email("invalid-email"),
				Password: "userPassword",
			},
			expectedStatusCode: 400,
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:    "INVALID_PARAMETER",
					Message: "invalid email format",
				},
			},
		},
		{
			name: "incorrect password",
			requestBody: dtos.LoginRequest{
				Email:    openapi_types.Email(email),
				Password: "wrongPassword",
			},
			expectedStatusCode: 401,
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:    "INVALID_CREDENTIALS",
					Message: "invalid email or password",
				},
			},
		},
		{
			name: "non-existent email",
			requestBody: dtos.LoginRequest{
				Email:    openapi_types.Email("nonexistent@example.com"),
				Password: "somePassword",
			},
			expectedStatusCode: 401,
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:    "INVALID_CREDENTIALS",
					Message: "invalid email or password",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			var requestBody bytes.Buffer
			err := json.NewEncoder(&requestBody).Encode(tc.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s/login", testDomain, baseURLV1Auth), &requestBody)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			if tc.expectedStatusCode == http.StatusOK {
				var authResp dtos.AuthResponse
				err = json.NewDecoder(resp.Body).Decode(&authResp)
				require.NoError(t, err)
				require.NotEmpty(t, authResp.AccessToken)

				// parse and validate the JWT access token
				validateToken(t, authResp.AccessToken, user.Id.String())

				// check the cookie for the refresh token is set
				cookies := resp.Cookies()
				var refreshTokenCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						refreshTokenCookie = cookie
						break
					}
				}

				assert.NotNil(t, refreshTokenCookie, "refresh_token cookie not found in response")
				assert.NotEmpty(t, refreshTokenCookie.Value, "refresh_token cookie value is empty")

				validateRefreshTokenCookie(t, refreshTokenCookie, user.Id.String())

				return
			}

			var errorResp dtos.ErrorEnvelope
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedError.Error.Code, errorResp.Error.Code)
			assert.Equal(t, tc.expectedError.Error.Message, errorResp.Error.Message)
		})
	}
}

func validateToken(t *testing.T, tokenString string, expectedUserID string) {
	// 1. Basic format check
	parts := strings.Split(tokenString, ".")
	require.Equal(t, 3, len(parts), "JWT should have 3 parts")

	// 2. Parse without verification first (to inspect claims)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	require.NoError(t, err, "should be able to parse token without verification")

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	// 3. Validate standard claims
	exp, ok := claims["exp"].(float64)
	require.True(t, ok, "exp claim should exist and be numeric")

	expTime := time.Unix(int64(exp), 0)
	assert.True(t, expTime.After(time.Now()), "token should not be expired")
	assert.True(t, expTime.Before(time.Now().Add(25*time.Hour)), "token expiry should be reasonable (< 24h)")

	iat, ok := claims["iat"].(float64)
	require.True(t, ok, "iat claim should exist")
	iatTime := time.Unix(int64(iat), 0)
	assert.True(t, iatTime.Before(time.Now().Add(1*time.Minute)), "iat should be recent")

	// 4. Validate custom claims
	sub, ok := claims["sub"].(string)
	require.True(t, ok, "sub claim should exist and be string")
	assert.Equal(t, expectedUserID, sub)

	// 5. Validate security claims
	alg := token.Header["alg"]
	assert.NotEqual(t, "none", alg, "algorithm should not be 'none'")
	assert.Contains(t, []string{"HS256", "RS256", "ES256"}, alg, "should use secure algorithm")
}

func validateRefreshTokenCookie(t *testing.T, cookie *http.Cookie, userID string) {
	// 1. Basic presence
	require.NotNil(t, cookie, "refresh_token cookie should be present")
	require.NotEmpty(t, cookie.Value)

	// 2. Security attributes
	assert.True(t, cookie.HttpOnly, "refresh token must be HttpOnly")
	assert.True(t, cookie.Secure, "refresh token must be Secure in production")
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite, "should use SameSite=Strict")

	// 3. Path and domain
	assert.Equal(t, "/api/v1/auth", cookie.Path, "should be scoped to auth endpoints")

	// 4. Expiration
	assert.True(t, cookie.Expires.After(time.Now()), "cookie should not be expired")
	assert.True(t, cookie.Expires.Before(time.Now().Add(31*24*time.Hour)),
		"refresh token should expire in reasonable time (~30 days)")

	validateToken(t, cookie.Value, userID)
}
