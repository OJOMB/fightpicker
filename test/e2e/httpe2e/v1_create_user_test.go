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
			}`, takenUsername, takenEmail),
			expectedStatusCode: http.StatusCreated,
			expectedResponse: dtos.UserResponse{
				Username:  takenUsername,
				Email:     new(types.Email(takenEmail)),
				FirstName: new("John"),
				LastName:  new("Doe"),
				Bio:       "Just a test user.",
				Location:  new("Testville"),
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
					Message:   "invalid email format",
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
					Code:      "CONFLICTING_RESOURCES",
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
					Code:      "CONFLICTING_RESOURCES",
					Message:   "username already taken",
					RequestId: "req-id",
				},
			},
		},
	}

	client := &http.Client{}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", testDomain, baseURLV1Users), strings.NewReader(tc.requestBody))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if resp.StatusCode == http.StatusCreated {
				// check the response body does not contain the password field
				bodyBytes, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var userResp dtos.UserResponse
				require.NoError(t, json.Unmarshal(bodyBytes, &userResp))

				assert.NotContains(t, strings.ToLower(string(bodyBytes)), "password")

				// check id is valid uuid
				_, err = uuid.FromString(userResp.Id.String())
				require.NoError(t, err)

				assert.Equal(t, tc.expectedResponse.Username, userResp.Username, "usernames should match")
				assert.Equal(t, tc.expectedResponse.Email, userResp.Email, "emails should match")
				assert.Equal(t, tc.expectedResponse.FirstName, userResp.FirstName, "first names should match")
				assert.Equal(t, tc.expectedResponse.LastName, userResp.LastName, "last names should match")
				assert.Equal(t, tc.expectedResponse.Dob, userResp.Dob, "dates of birth should match")
				assert.Equal(t, tc.expectedResponse.Bio, userResp.Bio, "bios should match")
				assert.Equal(t, tc.expectedResponse.Location, userResp.Location, "locations should match")
				assert.Equal(t, tc.expectedResponse.ProfilePicture, userResp.ProfilePicture, "profile pictures should match")

				assert.False(t, userResp.CreatedAt.IsZero(), "createdAt should not be zero")
				assert.False(t, userResp.UpdatedAt.IsZero(), "updatedAt should not be zero")

				return
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var errResp dtos.ErrorEnvelope
			require.NoError(t, json.Unmarshal(bodyBytes, &errResp))

			assert.Equal(t, tc.expectedError.Error.Code, errResp.Error.Code)
			assert.Equal(t, tc.expectedError.Error.Message, errResp.Error.Message)
		})
	}
}
