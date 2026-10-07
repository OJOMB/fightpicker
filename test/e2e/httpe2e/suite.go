package httpe2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
)

const (
	testDomain = "http://localhost:8080"

	adminEmail    = "admin@fightpicker.com"
	adminPassword = "chanko"

	baseURLV1Users = "/api/v1/users"
	baseURLV1Auth  = "/api/v1/auth"
)

func newRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func newRandomEmail() string {
	return newRandomString(10) + "@example.com"
}

func newRandomUsername() string {
	return newRandomString(8)
}

func createTestUser(t *testing.T, email, password string) dtos.UserResponse {
	bio := newRandomString(20)
	location := newRandomString(15)

	userRequest := dtos.UserCreateReq{
		Email:     openapi_types.Email(email),
		Password:  password,
		Username:  newRandomUsername(),
		Bio:       &bio,
		FirstName: newRandomString(5),
		LastName:  newRandomString(5),
		Location:  &location,
		Dob: openapi_types.Date{
			Time: time.Date(1990, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		Gender: dtos.UserCreateReqGenderOther,
	}

	var requestBody bytes.Buffer
	err := json.NewEncoder(&requestBody).Encode(userRequest)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", testDomain, baseURLV1Users), &requestBody)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", "user-creation-"+email)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var userResponse dtos.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	require.NoError(t, err)

	return userResponse
}

type testUser struct {
	Id       string
	email    string
	password string
}

func cleanupTestUsers(t *testing.T, testUsers ...testUser) {
	for _, testUser := range testUsers {
		// login as the test user to obtain an access token for deletion
		loginBody := bytes.NewBuffer(fmt.Appendf(nil, `{"email": "%s", "password": "%s"}`, testUser.email, testUser.password))
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s/login", testDomain, baseURLV1Auth), loginBody)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-Id", "login-for-test-user-deletion-"+testUser.Id)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResp dtos.AuthResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err)

		req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s/%s", testDomain, baseURLV1Users, testUser.Id), nil)
		require.NoError(t, err)
		accessToken := loginResp.AccessToken

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		resp, err = client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}
}
