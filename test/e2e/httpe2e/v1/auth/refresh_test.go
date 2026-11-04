package httpe2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	"github.com/OJOMB/fightpicker/test/e2e"
	"github.com/OJOMB/fightpicker/test/e2e/httpe2e"
)

// the service only treats reuse as suspicious once the grace window has elapsed (15s in the service)
// this just has to be greater than the config value in the local environment
const refreshReuseSuspiciousDelay = 6 * time.Second

func TestV1Refresh(t *testing.T) {
	email := e2e.NewRandomEmail()
	password := e2e.NewRandomString(12)
	user := httpe2e.CreateTestUserV1(t, email, password)

	t.Cleanup(func() {
		httpe2e.CleanupTestUsers(t, httpe2e.TestUser{
			Id:       user.Id.String(),
			Email:    email,
			Password: password,
		})
	})

	t.Run("successful refresh rotates the refresh token", func(t *testing.T) {
		oldToken := loginForRefreshToken(t, email, password)

		resp := postRefresh(t, oldToken)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var authResp dtos.AuthResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&authResp))
		require.NotEmpty(t, authResp.AccessToken)
		validateToken(t, authResp.AccessToken, user.Id.String())

		newCookie := findRefreshCookie(resp)
		require.NotNil(t, newCookie, "refresh_token cookie not found in response")
		assert.NotEmpty(t, newCookie.Value)
		assert.NotEqual(t, oldToken, newCookie.Value, "refresh token should be rotated")

		validateRefreshTokenCookie(t, newCookie, user.Id.String())
	})

	t.Run("rotated refresh token can be used again", func(t *testing.T) {
		oldToken := loginForRefreshToken(t, email, password)

		resp := postRefresh(t, oldToken)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		newCookie := findRefreshCookie(resp)
		require.NotNil(t, newCookie)

		resp2 := postRefresh(t, newCookie.Value)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		var authResp dtos.AuthResponse
		require.NoError(t, json.NewDecoder(resp2.Body).Decode(&authResp))
		assert.NotEmpty(t, authResp.AccessToken)
	})

	type testCase struct {
		name               string
		cookie             *http.Cookie
		expectedStatusCode int
		expectedError      dtos.ErrorEnvelope
	}

	testCases := []testCase{
		{
			name:               "missing refresh token cookie",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:    "MISSING_REQUIRED_PARAMETER",
					Message: "missing refresh token",
				},
			},
		},
		{
			name:               "empty refresh token cookie",
			cookie:             &http.Cookie{Name: "refresh_token", Value: ""},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:    "MISSING_REQUIRED_PARAMETER",
					Message: "refresh_token: missing required parameter",
				},
			},
		},
		{
			name:               "unknown refresh token",
			cookie:             &http.Cookie{Name: "refresh_token", Value: "not-a-real-token"},
			expectedStatusCode: http.StatusNotFound,
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:    "RESOURCE_NOT_FOUND",
					Message: "resource not found",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			req, err := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s%s/refresh", httpe2e.TestDomain, httpe2e.BaseURLV1Auth),
				nil,
			)
			require.NoError(t, err)

			req.Header.Set("X-Request-Id", "reqID")
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}

			resp, err := (&http.Client{}).Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			var errorResp dtos.ErrorEnvelope
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&errorResp))

			assert.Equal(t, tc.expectedError.Error.Code, errorResp.Error.Code)
			assert.Equal(t, tc.expectedError.Error.Message, errorResp.Error.Message)
		})
	}
}

func TestV1RefreshTokenReuse(t *testing.T) {
	email := e2e.NewRandomEmail()
	password := e2e.NewRandomString(12)
	user := httpe2e.CreateTestUserV1(t, email, password)

	t.Cleanup(func() {
		httpe2e.CleanupTestUsers(t, httpe2e.TestUser{
			Id:       user.Id.String(),
			Email:    email,
			Password: password,
		})
	})

	t.Run("immediate reuse of a rotated token is rejected as non-suspicious", func(t *testing.T) {
		oldToken := loginForRefreshToken(t, email, password)

		resp := postRefresh(t, oldToken)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		reuseResp := postRefresh(t, oldToken)
		defer reuseResp.Body.Close()
		require.Equal(t, http.StatusUnauthorized, reuseResp.StatusCode)

		var errorResp dtos.ErrorEnvelope
		require.NoError(t, json.NewDecoder(reuseResp.Body).Decode(&errorResp))
		assert.Equal(t, "INVALID_CREDENTIALS", errorResp.Error.Code)
		assert.Equal(t, "refresh token has been reused", errorResp.Error.Message)

		// a non-suspicious reuse must not revoke the family so the rotated token still works
		newCookie := findRefreshCookie(resp)
		require.NotNil(t, newCookie)

		validResp := postRefresh(t, newCookie.Value)
		defer validResp.Body.Close()
		assert.Equal(t, http.StatusOK, validResp.StatusCode)
	})

	t.Run("reuse of a rotated token after the grace window is flagged as suspicious", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping slow test in short mode")
		}

		oldToken := loginForRefreshToken(t, email, password)

		resp := postRefresh(t, oldToken)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		newCookie := findRefreshCookie(resp)
		require.NotNil(t, newCookie)

		time.Sleep(refreshReuseSuspiciousDelay)

		reuseResp := postRefresh(t, oldToken)
		defer reuseResp.Body.Close()
		require.Equal(t, http.StatusBadRequest, reuseResp.StatusCode)

		var errorResp dtos.ErrorEnvelope
		require.NoError(t, json.NewDecoder(reuseResp.Body).Decode(&errorResp))
		assert.Equal(t, "UNAUTHORIZED", errorResp.Error.Code)
		assert.Equal(t, "suspicious refresh token reuse", errorResp.Error.Message)

		// the failed refresh should clear the cookie
		cleared := findRefreshCookie(reuseResp)
		require.NotNil(t, cleared, "refresh_token cookie should be cleared")
		assert.Empty(t, cleared.Value)
		assert.Negative(t, cleared.MaxAge)

		// the whole token family is revoked so the legitimately rotated token must no longer work
		familyResp := postRefresh(t, newCookie.Value)
		defer familyResp.Body.Close()
		assert.NotEqual(t, http.StatusOK, familyResp.StatusCode, "rotated token should be revoked after suspicious reuse")
	})
}

// loginForRefreshToken logs in and returns the refresh token issued in the cookie
func loginForRefreshToken(t *testing.T, email, password string) string {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s/login", httpe2e.TestDomain, httpe2e.BaseURLV1Auth),
		strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)),
	)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", "login-for-refresh-"+email)

	resp, err := (&http.Client{}).Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	cookie := findRefreshCookie(resp)
	require.NotNil(t, cookie, "refresh_token cookie not found in login response")
	require.NotEmpty(t, cookie.Value)

	return cookie.Value
}

func postRefresh(t *testing.T, refreshToken string) *http.Response {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s/refresh", httpe2e.TestDomain, httpe2e.BaseURLV1Auth),
		nil,
	)
	require.NoError(t, err)

	req.Header.Set("X-Request-Id", "reqID")
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func findRefreshCookie(resp *http.Response) *http.Cookie {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "refresh_token" {
			return cookie
		}
	}

	return nil
}
