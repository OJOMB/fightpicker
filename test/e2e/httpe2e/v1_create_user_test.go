package httpe2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/fightpicker/internal/http/dtos"
)

func TestV1CreateUser(t *testing.T) {
	type testCase struct {
		name               string
		requestBody        string
		expectedStatusCode int
		expectedResponse   dtos.UserResponse
		expectedError      dtos.ErrorEnvelope
	}

	randEmail1 := newRandomEmail()
	randUsername1 := newRandomUsername()

	takenEmail := newRandomEmail()
	takenUsername := newRandomUsername()

	testCases := []testCase{
		{
			name: "successful user creation",
			requestBody: fmt.Sprintf(`{
				"username": "%s",
				"email": "%s",
				"first_name": "John",
				"last_name": "Doe",
				"bio": "Just a test user.",
				"location": "Testville",
				"dob": "1990-01-01",
				"password": "SecurePass123!"
			}`, randUsername1, randEmail1),
			expectedStatusCode: http.StatusCreated,
			expectedResponse: dtos.UserResponse{
				Username:  randUsername1,
				Email:     ptrOAPIEmail(types.Email(randEmail1)),
				FirstName: ptrString("John"),
				LastName:  ptrString("Doe"),
				Bio:       "Just a test user.",
				Location:  ptrString("Testville"),
				Dob: &types.Date{
					Time: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "missing required email field",
			requestBody: `{
				"username": "testuser2",
				"first_name": "Jane",
				"last_name": "Doe",
				"bio": "Another test user.",
				"location": "Exampletown",
				"dob": "1992-02-02",
				"password": "AnotherSecurePass123!"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   dtos.UserResponse{},
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:      "MISSING_REQUIRED_PARAMETER",
					Message:   "email: missing parameter",
					RequestId: "req-id",
				},
			},
		},
		{
			name: "missing required first_name field",
			requestBody: `{
				"username": "testuser2",
				"email": "testuser2@example.com",
				"last_name": "Doe",
				"bio": "Another test user.",
				"location": "Exampletown",
				"dob": "1992-02-02",
				"password": "AnotherSecurePass123!"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   dtos.UserResponse{},
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:      "MISSING_REQUIRED_PARAMETER",
					Message:   "first_name: missing parameter",
					RequestId: "req-id",
				},
			},
		},
		{
			name: "missing required last_name field",
			requestBody: `{
				"username": "testuser2",
				"email": "testuser2@example.com",
				"first_name": "Jane",
				"bio": "Another test user.",
				"location": "Exampletown",
				"dob": "1992-02-02",
				"password": "AnotherSecurePass123!"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   dtos.UserResponse{},
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:      "MISSING_REQUIRED_PARAMETER",
					Message:   "last_name: missing parameter",
					RequestId: "req-id",
				},
			},
		},
		{
			name: "missing required username field",
			requestBody: `{
				"email": "testuser2@example.com",
				"first_name": "Jane",
				"last_name": "Doe",
				"bio": "Another test user.",
				"location": "Exampletown",
				"dob": "1992-02-02",
				"password": "AnotherSecurePass123!"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   dtos.UserResponse{},
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:      "MISSING_REQUIRED_PARAMETER",
					Message:   "username: missing parameter",
					RequestId: "req-id",
				},
			},
		},
		{
			name: "missing required dob field",
			requestBody: `{
				"username": "testuser2",
				"email": "testuser2@example.com",
				"first_name": "Jane",
				"last_name": "Doe",
				"bio": "Another test user.",
				"location": "Exampletown",
				"password": "AnotherSecurePass123!"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   dtos.UserResponse{},
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:      "MISSING_REQUIRED_PARAMETER",
					Message:   "dob: missing parameter",
					RequestId: "req-id",
				},
			},
		},
		{
			name: "missing required password field",
			requestBody: `{
				"username": "testuser2",
				"email": "testuser2@example.com",
				"first_name": "Jane",
				"last_name": "Doe",
				"bio": "Another test user.",
				"location": "Exampletown",
				"dob": "1992-02-02"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   dtos.UserResponse{},
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:      "MISSING_REQUIRED_PARAMETER",
					Message:   "password: missing parameter",
					RequestId: "req-id",
				},
			},
		},
		{
			name: "invalid email format",
			requestBody: `{
				"username": "testuser3",
				"email": "invalid-email-format",
				"first_name": "Invalid",
				"last_name": "Email",
				"bio": "Testing invalid email.",
				"location": "Nowhere",
				"dob": "1995-03-03",
				"password": "InvalidEmailPass123!"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   dtos.UserResponse{},
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:      "INVALID_PARAMETER",
					Message:   "email: failed to pass regex validation",
					RequestId: "req-id",
				},
			},
		},
		{
			name: "duplicate email",
			requestBody: fmt.Sprintf(`{
				"username": "%s",
				"email": "%s",
				"first_name": "Duplicate",
				"last_name": "Email",
				"bio": "Testing duplicate email.",
				"location": "Somewhere",
				"dob": "1993-04-04",
				"password": "DuplicateEmailPass123!"
			}`, newRandomUsername(), takenEmail),
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   dtos.UserResponse{},
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:      "CONFLICT",
					Message:   "email already taken",
					RequestId: "req-id",
				},
			},
		},
		{
			name: "duplicate username",
			requestBody: fmt.Sprintf(`{
				"username": "%s",
				"email": "%s",
				"first_name": "Duplicate",
				"last_name": "Email",
				"bio": "Testing duplicate email.",
				"location": "Somewhere",
				"dob": "1993-04-04",
				"password": "DuplicateEmailPass123!"
			}`, takenUsername, newRandomEmail()),
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   dtos.UserResponse{},
			expectedError: dtos.ErrorEnvelope{
				Error: dtos.ErrorObject{
					Code:      "CONFLICT",
					Message:   "username already taken",
					RequestId: "req-id",
				},
			},
		},
	}

	// Setup: create a user to test duplicate email and username cases
	// we don't use the original success test case to avoid dependency on it passing and ordering issues with t.Run
	duplicateSetupRequestBody := fmt.Sprintf(`{
		"username": "%s",
		"email": "%s",
		"first_name": "Setup",
		"last_name": "User",
		"bio": "Setup user for duplicate tests.",
		"location": "Setupville
		"dob": "1991-01-01",
		"password": "SetupPass123!"
	}`, takenUsername, takenEmail)

	setupreq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", testDomain, baseURLV1Users), strings.NewReader(duplicateSetupRequestBody))
	require.NoError(t, err)

	setupreq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(setupreq)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", testDomain, baseURLV1Users), strings.NewReader(tc.requestBody))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Request-ID", "req-id")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			assert.Equal(t, "req-id", resp.Header.Get("X-Request-ID"))

			if resp.StatusCode == http.StatusCreated {
				// check the response body does not contain the password field
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				bodyString := string(bodyBytes)
				assert.NotContains(t, bodyString, "password")

				var userResp dtos.UserResponse
				require.NoError(t, json.Unmarshal(bodyBytes, &userResp))

				// check id is valid uuid
				_, err = uuid.FromString(userResp.Id.String())
				require.NoError(t, err)

				assert.Equal(t, tc.expectedResponse.Username, userResp.Username)
				assert.Equal(t, tc.expectedResponse.Email, userResp.Email)
				assert.Equal(t, tc.expectedResponse.FirstName, userResp.FirstName)
				assert.Equal(t, tc.expectedResponse.LastName, userResp.LastName)
				assert.Equal(t, tc.expectedResponse.Dob, userResp.Dob)
				assert.Equal(t, tc.expectedResponse.Bio, userResp.Bio)
				assert.Equal(t, tc.expectedResponse.Location, userResp.Location)
				assert.Equal(t, tc.expectedResponse.ProfilePicture, userResp.ProfilePicture)

				assert.False(t, userResp.CreatedAt.IsZero())
				assert.False(t, userResp.UpdatedAt.IsZero())

				return
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var errResp dtos.ErrorEnvelope
			err = json.Unmarshal(bodyBytes, &errResp)
			if err != nil {
				t.Logf("Failed to decode error response. Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))
				require.NoError(t, err, "Response body should be valid JSON")
			}

			assert.Equal(t, tc.expectedError.Error.Code, errResp.Error.Code)
			assert.Equal(t, tc.expectedError.Error.Message, errResp.Error.Message)
			assert.Equal(t, tc.expectedError.Error.RequestId, errResp.Error.RequestId)
		})
	}
}
