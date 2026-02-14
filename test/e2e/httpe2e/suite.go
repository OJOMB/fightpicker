package httpe2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/oapi-codegen/runtime/types"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	"github.com/OJOMB/fightpicker/pkg/id"
)

const (
	testDomain = "http://localhost:8080"

	baseURLV1Users = "/api/v1/users"
	baseURLV1Auth  = "/api/v1/auth"
)

func ptrString(s string) *string {
	return &s
}

func ptrOAPIEmail(email types.Email) *types.Email {
	return &email
}

func newRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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

func cleanupTestUser(t *testing.T, userID id.UUID7) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s/%s", testDomain, baseURLV1Users, userID.String()), nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}
