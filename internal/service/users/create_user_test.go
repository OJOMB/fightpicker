package users

import (
	"testing"
	"time"
)

func TestServiceCreateUser(t *testing.T) {
	type testCase struct {
		name         string
		inputUser    User
		expectError  bool
		expectedUser User
	}

	testCases := []testCase{
		{
			name: "Valid User Creation",
			inputUser: User{
				Email:        "johndoe@example.com",
				FirstName:    "John",
				LastName:     "Doe",
				Username:     "johndoe",
				DOB:          time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				PasswordHash: "securepassword",
			},
			expectError: false,
			expectedUser: User{
				Email:        "johndoe@example.com",
				FirstName:    "John",
				LastName:     "Doe",
				Username:     "johndoe",
				DOB:          time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
				PasswordHash: "i_got_hashed_securepassword",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
		})
	}
}
