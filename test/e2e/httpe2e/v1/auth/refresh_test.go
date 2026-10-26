package httpe2e

import (
	"fmt"
	"testing"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
	"github.com/OJOMB/fightpicker/test/e2e"
	"github.com/OJOMB/fightpicker/test/e2e/httpe2e"
)

func TestV1Refresh(t *testing.T) {
	// create a test user to refresh with
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

	type testCase struct {
		name               string
		requestBody        string
		expectedStatusCode int
		expectedError      dtos.ErrorEnvelope
	}

	testCases := []testCase{
		{
			name:               "example test case",
			requestBody:        `{"refresh_token": "someToken"}`,
			expectedStatusCode: 200,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {})
	}
}
