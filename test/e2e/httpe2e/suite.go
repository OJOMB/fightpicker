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
	TestDomain = "http://localhost:8080"

	AdminEmail    = "admin@fightpicker.com"
	AdminPassword = "chanko"

	BaseURLV1Users = "/api/v1/users"
	BaseURLV1Auth  = "/api/v1/auth"
)

func NewRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func NewRandomEmail() string {
	return NewRandomString(10) + "@example.com"
}

func NewRandomUsername() string {
	return NewRandomString(8)
}

func CreateTestUser(t *testing.T, email, password string) dtos.UserResponse {
	bio := NewRandomString(20)
	location := NewRandomString(15)

	userRequest := dtos.UserCreateReq{
		Email:     openapi_types.Email(email),
		Password:  password,
		Username:  NewRandomUsername(),
		Bio:       &bio,
		FirstName: NewRandomString(5),
		LastName:  NewRandomString(5),
		Location:  &location,
		Dob: openapi_types.Date{
			Time: time.Date(1990, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		Gender: dtos.UserCreateReqGenderOther,
	}

	var requestBody bytes.Buffer
	err := json.NewEncoder(&requestBody).Encode(userRequest)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", TestDomain, BaseURLV1Users), &requestBody)
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

type TestUser struct {
	Id       string
	Email    string
	Password string
}

func CleanupTestUsers(t *testing.T, testUsers ...TestUser) {
	for _, testUser := range testUsers {
		// login as the test user to obtain an access token for deletion
		loginBody := bytes.NewBuffer(fmt.Appendf(nil, `{"email": "%s", "password": "%s"}`, testUser.Email, testUser.Password))
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s/login", TestDomain, BaseURLV1Auth), loginBody)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-Id", "login-for-test-user-deletion-"+testUser.Email)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResp dtos.AuthResponse
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err)

		req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s/%s", TestDomain, BaseURLV1Users, testUser.Id), nil)
		require.NoError(t, err)
		accessToken := loginResp.AccessToken

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		resp, err = client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}
}
